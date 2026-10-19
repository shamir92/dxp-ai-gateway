//go:build manual

// Operator tool: decrypt stored external-platform credentials.
//
// Not part of the normal build, so it can never ship in the server binary:
//
//	cd lp/backend
//	go run -tags manual ./cmd/decrypt -email someone@example.com
//	go run -tags manual ./cmd/decrypt -email someone@example.com -site tpai -what password
//	go run -tags manual ./cmd/decrypt -email someone@example.com -what token
//
// It loads .env (next to this module, then the repo-relative path) the same way
// the server does, so an exported shell environment is not required. Real
// environment variables always win over the file.
//
// Decrypted values are printed to STDOUT. They are live credentials — do not
// paste the output into tickets, logs, or chat.

package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/credentials"
	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/config"
)

func main() {
	email := flag.String("email", "", "gateway user email (required)")
	siteName := flag.String("site", "", "tpai | apai (default: all sites with a row)")
	what := flag.String("what", "all", "password | token | all")
	flag.Parse()

	loadDotEnv()
	if err := run(*email, *siteName, *what); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// loadDotEnv mirrors the server's lookup order. A missing file is fine; real
// environment variables are never overwritten.
func loadDotEnv() {
	paths := []string{".env", "lp/backend/.env"}
	for _, dir := range []string{".", ".."} {
		paths = append(paths, dir+"/lp/backend/.env")
	}
	for _, path := range paths {
		err := config.LoadDotEnv(path)
		if err == nil {
			fmt.Fprintf(os.Stderr, "loaded %s\n", path)
			return
		}
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: could not read %s: %v\n", path, err)
		}
	}
	fmt.Fprintln(os.Stderr, "warning: no .env found; relying on the shell environment")
}

func run(email, siteName, what string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("-email is required")
	}
	switch what {
	case "password", "token", "all":
	default:
		return fmt.Errorf("-what must be password, token, or all")
	}

	keyText := strings.TrimSpace(os.Getenv("EXTERNAL_PASSWORD_ENCRYPTION_KEY"))
	if keyText == "" {
		return fmt.Errorf("EXTERNAL_PASSWORD_ENCRYPTION_KEY is not set; checked the shell environment and .env")
	}
	key, err := base64.StdEncoding.DecodeString(keyText)
	if err != nil {
		return fmt.Errorf("EXTERNAL_PASSWORD_ENCRYPTION_KEY is not valid base64: %w", err)
	}
	cipher, err := credentials.NewCipher(key)
	if err != nil {
		return err
	}

	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set; checked the shell environment and .env")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	query := `
		SELECT e.site, e.username, e.password_ciphertext, e.password_nonce, e.encryption_key_version,
		       e.access_token_ciphertext, e.access_token_nonce, e.access_token_key_version,
		       e.access_token_status
		FROM external_site_accounts e
		JOIN users u ON u.id = e.user_id
		WHERE u.email = $1`
	args := []any{email}
	if siteName != "" {
		site := entity.ExternalSite(strings.ToLower(strings.TrimSpace(siteName)))
		if site != entity.ExternalSiteTPAI && site != entity.ExternalSiteAPAI {
			return fmt.Errorf("-site must be tpai or apai")
		}
		query += ` AND e.site = $2`
		args = append(args, site)
	}
	query += ` ORDER BY e.site`

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	found := 0
	for rows.Next() {
		var (
			site        entity.ExternalSite
			username    string
			pwCT, pwN   []byte
			pwVersion   string
			atCT, atN   []byte
			atVersion   *string
			atStatus    string
			atVersionOK string
		)
		if err := rows.Scan(&site, &username, &pwCT, &pwN, &pwVersion,
			&atCT, &atN, &atVersion, &atStatus); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		if atVersion != nil {
			atVersionOK = *atVersion
		}
		found++

		fmt.Printf("== %s ==\n", site)
		fmt.Printf("   remote username : %s\n", username)

		if what == "password" || what == "all" {
			pw, err := cipher.Decrypt(siteUserID(ctx, pool, email), site, pwCT, pwN, pwVersion)
			if err != nil {
				fmt.Printf("   password        : <cannot decrypt: %v>\n", err)
			} else {
				fmt.Printf("   password        : %s\n", pw)
			}
		}
		if what == "token" || what == "all" {
			if len(atCT) == 0 {
				fmt.Printf("   access token    : <none stored (status=%s)>\n", atStatus)
			} else {
				tok, err := cipher.DecryptAccessToken(siteUserID(ctx, pool, email), site, atCT, atN, atVersionOK)
				if err != nil {
					fmt.Printf("   access token    : <cannot decrypt: %v>\n", err)
				} else {
					fmt.Printf("   access token    : %s\n", tok)
				}
			}
		}
		fmt.Println()
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if found == 0 {
		return fmt.Errorf("no external_site_accounts row for %s", email)
	}
	return nil
}

// siteUserID resolves the user id, which is part of the cipher's associated data:
// a ciphertext only opens under the exact user and site it was sealed for.
func siteUserID(ctx context.Context, pool *pgxpool.Pool, email string) (id [16]byte) {
	if err := pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return id
		}
	}
	return id
}

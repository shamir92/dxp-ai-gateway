package config

import (
	"strings"
	"testing"
)

// The whole point of DatabaseTarget is that it is safe to log, so the password
// must never survive it. A regression here leaks a database credential into the
// log pipeline.
func TestDatabaseTargetNeverLeaksCredentials(t *testing.T) {
	const password = "sup3r-secret-pw"
	for _, tc := range []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "pooler with port and params",
			dsn:  "postgresql://postgres.abc:" + password + "@aws-0-ap-southeast-2.pooler.supabase.com:5432/postgres?sslmode=require",
			want: "aws-0-ap-southeast-2.pooler.supabase.com:5432/postgres",
		},
		{
			name: "local default",
			dsn:  "postgres://dxp:" + password + "@localhost:5432/dxp_ai_gateway?sslmode=disable",
			want: "localhost:5432/dxp_ai_gateway",
		},
		{
			name: "no port",
			dsn:  "postgres://user:" + password + "@db.example.com/mydb",
			want: "db.example.com/mydb",
		},
		{
			name: "key=value form cannot be parsed as a URL",
			dsn:  "host=localhost port=5432 dbname=dxp password=" + password,
			want: "unknown",
		},
		{
			name: "not a dsn",
			dsn:  "nonsense",
			want: "unknown",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := DatabaseTarget(tc.dsn)
			if got != tc.want {
				t.Errorf("DatabaseTarget() = %q, want %q", got, tc.want)
			}
			if strings.Contains(got, password) {
				t.Fatalf("DatabaseTarget() leaked the password: %q", got)
			}
		})
	}
}

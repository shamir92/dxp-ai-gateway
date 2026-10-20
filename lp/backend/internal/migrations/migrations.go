// Package migrations embeds the SQL migration files so the server can migrate
// without the files being present on disk at runtime.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS

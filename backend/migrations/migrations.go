// Package migrations embeds the SQL migration files into the binary so the
// API container can apply them on startup (and cmd/migrate can run them)
// without shipping loose files.
package migrations

import "embed"

//go:embed *.sql
var files embed.FS

// FS exposes the embedded migration files to golang-migrate's iofs source.
func FS() embed.FS { return files }

package sqlite

import _ "embed"

//go:generate -command sqlc go run github.com/kyleconroy/sqlc/cmd/sqlc@v1.17.2
//go:generate sqlc generate

//go:embed schema.sql
var Schema string // expose the schema for connecting

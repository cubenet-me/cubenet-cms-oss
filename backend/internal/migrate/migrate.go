package migrate

import "embed"

//go:embed migrations/*.sql
var FS embed.FS

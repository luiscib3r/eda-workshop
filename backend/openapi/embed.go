package openapi

import (
	"embed"
)

//go:embed all:*
var OpenApiFS embed.FS

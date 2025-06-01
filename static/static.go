package static

import "embed"

//go:embed css js img
var StaticFiles embed.FS

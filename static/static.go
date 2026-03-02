package static

import "embed"

//go:embed all:js all:css
var FS embed.FS

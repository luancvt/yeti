package static

import "embed"

//go:embed all:js all:css all:fonts
var FS embed.FS

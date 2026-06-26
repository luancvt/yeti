package static

import "embed"

//go:embed all:js all:css all:fonts favicon.png
var FS embed.FS

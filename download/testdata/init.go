package testdata

import "embed"

//go:embed *.zip *.txt
var FS embed.FS

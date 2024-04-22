package ui

import (
	"embed"
	"io/fs"
)

//go:embed "html" "static"
var Files embed.FS
var StaticFS, _ = fs.Sub(Files, "static")

package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:ui
var frontendDist embed.FS

// StaticFS returns the embedded frontend file system
func StaticFS() (http.FileSystem, error) {
	// Root of the embedded FS is the 'ui' directory.
	subFS, err := fs.Sub(frontendDist, "ui")
	if err != nil {
		return nil, err
	}
	return http.FS(subFS), nil
}

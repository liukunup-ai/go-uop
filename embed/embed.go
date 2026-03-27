package embed

import (
	"embed"
	"io/fs"
	"net/http"
)

var ConsoleFS embed.FS

func GetFS() (fs.FS, error) {
	sub, err := fs.Sub(ConsoleFS, "_out")
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func HTTPFileSystem() http.FileSystem {
	f, err := GetFS()
	if err != nil {
		return http.Dir(".")
	}
	return http.FS(f)
}

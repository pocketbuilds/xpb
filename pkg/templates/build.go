package templates

import (
	"embed"
)

//go:embed build
var buildFS embed.FS

func CopyBuildTemplate(destDir string, data any) error {
	return copyTemplate(buildFS, "build", destDir, data)
}

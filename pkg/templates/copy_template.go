package templates

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func copyTemplate(srcFs fs.FS, rootDir string, destDir string, data any) error {
	return fs.WalkDir(srcFs, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		var destPath = path
		destPath = strings.TrimPrefix(destPath, rootDir)
		destPath = strings.TrimSuffix(destPath, ".tmpl")
		destPath = filepath.Join(destDir, destPath)

		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		if filepath.Ext(path) == ".tmpl" {
			tmpl, err := template.ParseFS(srcFs, path)
			if err != nil {
				return err
			}
			return tmpl.Execute(destFile, data)
		} else {
			srcFile, err := srcFs.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			_, err = io.Copy(destFile, srcFile)
			return err
		}
	})
}

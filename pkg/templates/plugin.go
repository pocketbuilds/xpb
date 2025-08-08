package templates

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed plugin
var pluginFS embed.FS

type PluginTemplateData struct {
	Name string
}

func GeneratePluginDir(destDir string, data PluginTemplateData) error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go toolchain missing: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "xpb_plugin")
	if err != nil {
		return fmt.Errorf("error making tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	err = copyTemplate(pluginFS, "plugin", tmpDir, data)
	if err != nil {
		return fmt.Errorf("error copying template: %w", err)
	}

	var cmd *exec.Cmd

	cmd = exec.Command("go", "mod", "init", data.Name)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running go mod init: %w", err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running go mod tidy: %w", err)
	}

	return filepath.WalkDir(tmpDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destPath := filepath.Join(destDir, strings.TrimPrefix(path, tmpDir))

		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		if err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

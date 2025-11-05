package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pocketbuilds/xpb/pkg/module"
	"github.com/pocketbuilds/xpb/pkg/templates"
)

func (b *Builder) Build() (r io.ReadCloser, err error) {

	const binFileName = "pocketbase"

	steps := []func() error{
		b.printGoInfo,
		b.copyBuildTemplate,
		b.runGoModInit,
		b.runGoGetForAllModules,
		b.handleModuleReplacements,
		b.addOptimizationLdFlags,
		b.addVersionLdFlags,
		b.runGoBuild(binFileName),
	}

	for _, runStep := range steps {
		if err := runStep(); err != nil {
			return nil, err
		}
	}

	return b.buildResult(filepath.Join(b.dir, binFileName))
}

func (b *Builder) printGoInfo() error {
	cmd := b.newCommand("go", "env", "GOPATH")
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *Builder) copyBuildTemplate() error {
	return templates.CopyBuildTemplate(b.dir, b)
}

func (b *Builder) runGoModInit() error {
	cmd := b.newCommand("go", "mod", "init", "pocketbase")
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *Builder) runGoGetForAllModules() error {

	cmd := b.newCommand("go", "get", "-v", b.Pocketbase.String())
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Including pocketbase in all go get commands to ensure go toolkit does not
	//   silently bump the pocketbase module version. See xcaddy issue below:
	// https://github.com/caddyserver/xcaddy/issues/54

	cmd = b.newCommand("go", "get", "-v", b.Pocketbase.String(), b.Xpb.String())
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, module := range b.Plugins {
		cmd = b.newCommand("go", "get", "-v", b.Pocketbase.String(), b.Xpb.String(), module.String())
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// https://github.com/caddyserver/xcaddy/pull/92
	cmd = b.newCommand("go", "get")
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (b *Builder) handleModuleReplacements() error {
	allModules := append([]*module.Module{b.Pocketbase, b.Xpb}, b.Plugins...)
	for _, module := range allModules {
		if module.Replacement == "" {
			continue
		}
		replacement, err := filepath.Abs(module.Replacement)
		if err != nil {
			return err
		}
		cmd := b.newCommand(
			"go", "mod", "edit",
			"-replace", module.Module+"="+replacement,
		)
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) addOptimizationLdFlags() error {
	b.LdFlags = append(b.LdFlags,
		"-s", "-w",
	)
	return nil
}

func (b *Builder) addVersionLdFlags() error {
	cmd := b.newCommand("go", "list", "-m", "all")
	cmd.Stdout = nil
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	allModules := append([]*module.Module{b.Pocketbase, b.Xpb}, b.Plugins...)
	for _, module := range allModules {
		re, err := regexp.Compile(module.Module + " (.+)+")
		if err != nil {
			return err
		}
		match := re.FindStringSubmatch(string(output))
		if match == nil {
			return nil
		}
		b.LdFlags = append(b.LdFlags,
			fmt.Sprintf("-X '%s.version=%s'", module.Module, match[1]),
		)
	}
	return nil
}

func (b *Builder) runGoBuild(filename string) func() error {
	return func() error {
		args := []string{
			"build",
			"-o", filename,
		}
		if len(b.Tags) != 0 {
			args = append(args,
				"-tags", strings.Join(b.Tags, ","),
			)
		}
		args = append(args,
			"-ldflags", strings.Join(b.LdFlags, " "),
		)
		cmd := b.newCommand("go", args...)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", b.Os),
			fmt.Sprintf("GOARCH=%s", b.Arch),
		)
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
}

func (b *Builder) buildResult(binFilePath string) (io.ReadCloser, error) {
	binFile, err := os.Open(binFilePath)
	if err != nil {
		return nil, err
	}

	return &buildReadCloser{
		file:  binFile,
		dir:   b.dir,
		rmDir: b.rmDir,
	}, nil
}

type buildReadCloser struct {
	file  *os.File
	dir   string
	rmDir bool
}

func (brc *buildReadCloser) Read(p []byte) (int, error) {
	return brc.file.Read(p)
}

func (brc *buildReadCloser) Close() error {
	if brc.rmDir {
		defer os.RemoveAll(brc.dir)
	}
	return brc.file.Close()
}

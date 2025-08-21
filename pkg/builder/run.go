package builder

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pocketbuilds/xpb/pkg/module"
	"github.com/pocketbuilds/xpb/pkg/templates"
)

func (b *Builder) Build() (r io.ReadCloser, err error) {

	const binFileName = "pocketbase"

	steps := []func() error{
		b.copyBuildTemplate,
		b.runGoModInit,
		b.handleGoModReplacements,
		b.runGoModTidy,
		b.runGoGetPocketbaseAtSpecifiedVersion,
		b.runGoGetXpbAtSpecifiedVersion,
		b.addPocketbaseVersionLdFlag,
		b.addXpbVersionLdFlag,
		b.addPluginVersionsLdFlag,
		b.runGoBuild(binFileName),
	}

	for _, runStep := range steps {
		if err := runStep(); err != nil {
			return nil, err
		}
	}

	return b.buildResult(path.Join(b.dir, binFileName))
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

func (b *Builder) handleGoModReplacements() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, module := range b.Plugins {
		if module.Replacement != "" {
			replacement := path.Join(wd, module.Replacement)
			cmd := b.newCommand(
				"go", "mod", "edit",
				"-replace", module.Module+"="+replacement,
			)
			fmt.Fprintf(b.stdout, "%s\n", cmd)
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Builder) runGoModTidy() error {
	cmd := b.newCommand("go", "mod", "tidy", "-v")
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *Builder) runGoGetPocketbaseAtSpecifiedVersion() error {
	cmd := b.newCommand("go", "get", "-v", b.Pocketbase.Module+"@"+b.Pocketbase.Version)
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	if b.Pocketbase.Replacement != "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		replacement := path.Join(wd, b.Pocketbase.Replacement)
		cmd := b.newCommand(
			"go", "mod", "edit",
			"-replace", b.Pocketbase.Module+"="+replacement,
		)
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) runGoGetXpbAtSpecifiedVersion() error {
	cmd := b.newCommand("go", "get", "-v", b.Xpb.Module+"@"+b.Xpb.Version)
	fmt.Fprintf(b.stdout, "%s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	if b.Xpb.Replacement != "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		replacement := path.Join(wd, b.Xpb.Replacement)
		cmd := b.newCommand(
			"go", "mod", "edit",
			"-replace", b.Xpb.Module+"="+replacement,
		)
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		cmd := b.newCommand("go", "list", "-m", "all")
		fmt.Fprintf(b.stdout, "%s\n", cmd)
		cmd.Stdout = nil
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		re := regexp.MustCompile(module.XpbModule + " (.+)")
		match := re.FindStringSubmatch(string(output))
		if len(match) > 1 {
			b.Xpb.Version = match[1]
		}
	}
	return nil
}

func (b *Builder) addPocketbaseVersionLdFlag() error {
	cmd := b.newCommand("go", "list", "-m", "all")
	cmd.Stdout = nil
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	re, err := regexp.Compile(b.Pocketbase.Module + " (.+)")
	if err != nil {
		return err
	}
	match := re.FindStringSubmatch(string(output))
	if match == nil {
		return nil
	}
	b.LdFlags = append(b.LdFlags,
		fmt.Sprintf("-X '%s.version=%s'", b.Pocketbase.Module, match[1]),
	)
	return nil
}

func (b *Builder) addXpbVersionLdFlag() error {
	cmd := b.newCommand("go", "list", "-m", "all")
	cmd.Stdout = nil
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	re, err := regexp.Compile(b.Xpb.Module + " (.+)")
	if err != nil {
		return err
	}
	match := re.FindStringSubmatch(string(output))
	if match == nil {
		return nil
	}
	b.LdFlags = append(b.LdFlags,
		fmt.Sprintf("-X '%s.version=%s'", b.Xpb.Module, match[1]),
	)
	return nil
}

func (b *Builder) addPluginVersionsLdFlag() error {
	cmd := b.newCommand("go", "list", "-m", "all")
	cmd.Stdout = nil
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	for _, p := range b.Plugins {
		re, err := regexp.Compile(p.Module + " (.+)")
		if err != nil {
			return err
		}
		match := re.FindStringSubmatch(string(output))
		if match == nil {
			continue
		}
		b.LdFlags = append(b.LdFlags,
			fmt.Sprintf("-X '%s.version=%s'", p.Module, match[1]),
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

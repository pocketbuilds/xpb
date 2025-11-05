package builder

import (
	"io"
	"os"
	"path/filepath"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbuilds/xpb/pkg/module"
	"github.com/pocketbuilds/xpb/pkg/rules"
)

type Builder struct {
	Arch    string   `json:"arch"`
	Os      string   `json:"os"`
	Tags    []string `json:"tags"`
	LdFlags []string `json:"ldflags"`

	Plugins    []*module.Module `json:"plugins"`
	Pocketbase *module.Module   `json:"pocketbase"`
	Xpb        *module.Module   `json:"xpb"`

	dir    string
	rmDir  bool
	stdout io.Writer
	stderr io.Writer
}

type BuilderOption func(b *Builder) error

func NewBuilder(opts ...BuilderOption) (*Builder, error) {
	var b = Builder{
		Xpb: &module.Module{
			Module: module.XpbModule,
		},
		Pocketbase: &module.Module{
			Module: module.PocketbaseModule,
		},
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, opt := range opts {
		if err := opt(&b); err != nil {
			return nil, err
		}
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	if b.dir == "" {
		var err error
		b.dir, err = os.MkdirTemp("", "pocketbase")
		if err != nil {
			return nil, err
		}
		b.rmDir = true
	}
	return &b, nil
}

func (b Builder) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Arch,
			validation.Required,
			rules.IsArch(),
		),
		validation.Field(&b.Os,
			validation.Required,
			rules.IsOs(),
		),
		validation.Field(&b.Plugins),
		validation.Field(&b.Xpb,
			validation.Required,
		),
		validation.Field(&b.Pocketbase,
			validation.Required,
		),
	)
}

func (b *Builder) BuildToFile(path string) error {
	absFilePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rc, err := b.Build()
	if err != nil {
		return err
	}
	defer rc.Close()

	file, err := os.OpenFile(absFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, rc)
	return err
}

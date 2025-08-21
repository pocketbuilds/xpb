package builder

import (
	"encoding/json"
	"io"

	"github.com/pocketbuilds/xpb/pkg/module"
)

func WithArch(arch string) BuilderOption {
	return func(b *Builder) error {
		b.Arch = arch
		return nil
	}
}

func WithOS(os string) BuilderOption {
	return func(b *Builder) error {
		b.Os = os
		return nil
	}
}

func WithPbVersion(version string) BuilderOption {
	return func(b *Builder) error {
		b.Pocketbase.Version = version
		return nil
	}
}

func WithStdoutWriter(w io.Writer) BuilderOption {
	return func(b *Builder) error {
		b.stdout = w
		return nil
	}
}

func WithStderrWriter(w io.Writer) BuilderOption {
	return func(b *Builder) error {
		b.stderr = w
		return nil
	}
}

func WithOutputWriter(w io.Writer) BuilderOption {
	return func(b *Builder) error {
		b.stdout = w
		b.stderr = w
		return nil
	}
}

func WithPlugins(plugins ...*module.Module) BuilderOption {
	return func(b *Builder) error {
		for _, p := range plugins {
			switch {
			case p.IsPocketbase():
				b.Pocketbase = p
			case p.IsXpb():
				b.Xpb = p
			default:
				b.Plugins = append(b.Plugins, p)
			}
		}
		return nil
	}
}

func WithBuildDir(dir string) BuilderOption {
	return func(b *Builder) error {
		b.dir = dir
		b.rmDir = false
		return nil
	}
}

func WithNewPlugin(opts ...module.ModuleOption) BuilderOption {
	return func(b *Builder) error {
		m, err := module.NewModule(opts...)
		if err != nil {
			return err
		}
		return WithPlugins(m)(b)
	}
}

func WithTags(tags ...string) BuilderOption {
	return func(b *Builder) error {
		b.Tags = append(b.Tags, tags...)
		return nil
	}
}

func WithLdflags(ldflags ...string) BuilderOption {
	return func(b *Builder) error {
		b.LdFlags = append(b.LdFlags, ldflags...)
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// JSON
///////////////////////////////////////////////////////////////////////////////

func FromJsonReader(r io.Reader) BuilderOption {
	return func(b *Builder) error {
		return json.NewDecoder(r).Decode(&b)
	}
}

func FromJsonBytes(data []byte) BuilderOption {
	return func(b *Builder) error {
		return json.Unmarshal(data, b)
	}
}

func FromJsonString(data string) BuilderOption {
	return FromJsonBytes([]byte(data))
}

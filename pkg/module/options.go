package module

import (
	"encoding/json"
	"io"
	"regexp"
)

func WithModule(module string) ModuleOption {
	return func(m *Module) error {
		m.Module = module
		return nil
	}
}

func WithVersion(version string) ModuleOption {
	return func(m *Module) error {
		m.Version = version
		return nil
	}
}

func WithReplacement(replacement string) ModuleOption {
	return func(m *Module) error {
		m.Replacement = replacement
		return nil
	}
}

func FromJsonReader(r io.Reader) ModuleOption {
	return func(m *Module) error {
		return json.NewDecoder(r).Decode(m)
	}
}

func FromJsonBytes(data []byte) ModuleOption {
	return func(m *Module) error {
		return json.Unmarshal(data, m)
	}
}

func FromJsonString(data string) ModuleOption {
	return FromJsonBytes([]byte(data))
}

var reCliArg = regexp.MustCompile(`^([^\s@=]+)@?([^\s@=]*)=?([^\s@=]*)$`)

func FromCliArg(arg string) ModuleOption {
	return func(m *Module) error {
		match := reCliArg.FindStringSubmatch(arg)
		if len(match) > 3 {
			m.Module = match[1]
			m.Version = match[2]
			m.Replacement = match[3]
		}
		return nil
	}
}

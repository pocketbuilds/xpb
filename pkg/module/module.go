package module

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Module struct {
	Module      string `json:"module"`
	Version     string `json:"version"`
	Replacement string `json:"replacement"`
}

type ModuleOption func(m *Module) error

func NewModule(opts ...ModuleOption) (*Module, error) {
	var module Module
	for _, opt := range opts {
		if err := opt(&module); err != nil {
			return nil, err
		}
	}
	return &module, nil
}

func (m Module) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Module,
			validation.Required,
		),
	)
}

func (m Module) String() string {
	str := m.Module
	if m.Version != "" && m.Version != "latest" {
		str += "@" + m.Version
	}
	return str
}

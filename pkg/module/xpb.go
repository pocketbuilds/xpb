package module

const XpbModule = "github.com/pocketbuilds/xpb"

func WithXpbModule() ModuleOption {
	return func(m *Module) error {
		m.Module = XpbModule
		return nil
	}
}

func (m *Module) IsXpb() bool {
	return m.Module == XpbModule
}

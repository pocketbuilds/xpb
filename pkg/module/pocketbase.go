package module

const PocketbaseModule = "ganigeorgiev.com/pocketbase/pocketbase"

func WithPocketbaseModule() ModuleOption {
	return func(m *Module) error {
		m.Module = PocketbaseModule
		return nil
	}
}

func (m *Module) IsPocketbase() bool {
	return m.Module == PocketbaseModule
}

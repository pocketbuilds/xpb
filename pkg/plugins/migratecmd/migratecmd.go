package migratecmd

import (
	"os"

	"github.com/PocketBuilds/xpb"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func init() {
	xpb.Register(&Plugin{
		Dir:          "",
		Automigrate:  true,
		TemplateLang: migratecmd.TemplateLangJS,
	})
}

type Plugin struct {
	// Dir specifies the directory with the user defined migrations.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" (for js)
	// or "pb_data/../migrations" (for go) directory.
	Dir string

	// Automigrate specifies whether to enable automigrations.
	Automigrate bool

	// TemplateLang specifies the template language to use when
	// generating migrations - js or go (default).
	TemplateLang string
}

func (p *Plugin) Name() string {
	return "migratecmd"
}

func (p *Plugin) Version() string {
	return xpb.Version()
}

func (p *Plugin) Description() string {
	return "The built-in pocketbase migratecmd plugin included as an xpb plugin."
}

func (p *Plugin) Init(app core.App) error {

	if app, ok := app.(*pocketbase.PocketBase); ok {

		app.RootCmd.PersistentFlags().StringVar(
			&p.Dir,
			"migrationsDir",
			p.Dir,
			"the directory with the user defined migrations",
		)

		app.RootCmd.PersistentFlags().BoolVar(
			&p.Automigrate,
			"automigrate",
			p.Automigrate,
			"enable/disable auto migrations",
		)

		app.RootCmd.ParseFlags(os.Args[1:])

		migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
			Automigrate:  p.Automigrate,
			Dir:          p.Dir,
			TemplateLang: p.TemplateLang,
		})
	}

	return nil
}

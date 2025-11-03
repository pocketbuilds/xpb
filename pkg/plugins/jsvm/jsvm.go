package jsvm

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbuilds/xpb"
)

func init() {
	xpb.Register(&Plugin{
		HooksDir:      "",
		HooksWatch:    true,
		HooksPoolSize: 25,
	})
}

type Plugin struct {

	// HooksWatch enables auto app restarts when a JS app hook file changes.
	//
	// Note that currently the application cannot be automatically restarted on Windows
	// because the restart process relies on execve.
	HooksWatch bool

	// HooksDir specifies the JS app hooks directory.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_hooks" directory.
	HooksDir string

	// HooksFilesPattern specifies a regular expression pattern that
	// identify which file to load by the hook vm(s).
	//
	// If not set it fallbacks to `^.*(\.pb\.js|\.pb\.ts)$`, aka. any
	// HookdsDir file ending in ".pb.js" or ".pb.ts" (the last one is to enforce IDE linters).
	HooksFilesPattern string

	// HooksPoolSize specifies how many goja.Runtime instances to prewarm
	// and keep for the JS app hooks gorotines execution.
	//
	// Zero or negative value means that it will create a new goja.Runtime
	// on every fired goroutine.
	HooksPoolSize int

	// MigrationsDir specifies the JS migrations directory.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" directory.
	MigrationsDir string

	// If not set it fallbacks to `^.*(\.js|\.ts)$`, aka. any MigrationDir file
	// ending in ".js" or ".ts" (the last one is to enforce IDE linters).
	MigrationsFilesPattern string

	// TypesDir specifies the directory where to store the embedded
	// TypeScript declarations file.
	//
	// If not set it fallbacks to "pb_data".
	//
	// Note: Avoid using the same directory as the HooksDir when HooksWatch is enabled
	// to prevent unnecessary app restarts when the types file is initially created.
	TypesDir string
}

// Name implements [xpb.Plugin.Name] interface method.
func (p *Plugin) Name() string {
	return "jsvm"
}

// Version implements [xpb.Plugin.Version] interface method.
func (p *Plugin) Version() string {
	return xpb.Version()
}

// Description implements [xpb.Plugin.Description] interface method.
func (p *Plugin) Description() string {
	return "The built-in pocketbase jsvm plugin included as an xpb plugin."
}

// Init implements [xpb.Plugin.Init] interface method.
func (p *Plugin) Init(app core.App) error {

	if app, ok := app.(*pocketbase.PocketBase); ok {
		app.RootCmd.PersistentFlags().StringVar(
			&p.HooksDir,
			"hooksDir",
			p.HooksDir,
			"the directory with the JS app hooks",
		)

		app.RootCmd.PersistentFlags().BoolVar(
			&p.HooksWatch,
			"hooksWatch",
			p.HooksWatch,
			"auto restart the app on pb_hooks file change",
		)

		app.RootCmd.PersistentFlags().IntVar(
			&p.HooksPoolSize,
			"hooksPool",
			p.HooksPoolSize,
			"the total prewarm goja.Runtime instances for the JS app hooks execution",
		)

		app.RootCmd.ParseFlags(os.Args[1:])

	}

	jsvm.MustRegister(app, jsvm.Config{
		HooksWatch:             p.HooksWatch,
		HooksDir:               p.HooksDir,
		HooksFilesPattern:      p.HooksFilesPattern,
		HooksPoolSize:          p.HooksPoolSize,
		MigrationsDir:          p.MigrationsDir,
		MigrationsFilesPattern: p.MigrationsFilesPattern,
		TypesDir:               p.TypesDir,
		OnInit:                 p.jsvmInitHandler(app),
	})

	return nil
}

type JsvmPlugin interface {
	OnJsvmInit(app core.App, vm *goja.Runtime)
}

func (p *Plugin) jsvmInitHandler(app core.App) func(vm *goja.Runtime) {
	return func(vm *goja.Runtime) {
		for _, plugin := range xpb.GetPlugins() {
			if jsvmPlugin, ok := plugin.(JsvmPlugin); ok {
				fmt.Println(plugin.Name())
				jsvmPlugin.OnJsvmInit(app, vm)
			}
		}
	}
}

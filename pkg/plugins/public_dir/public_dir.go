package public_dir

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbuilds/xpb"
)

func init() {
	xpb.Register(&Plugin{
		Dir:           defaultPublicDir(),
		IndexFallback: true,
	})
}

type Plugin struct {
	// the directory to serve static files
	Dir string
	// fallback the request to index.html on missing static path (eg. when pretty urls are used with SPA)
	IndexFallback bool
}

func (p *Plugin) Name() string {
	return "public_dir"
}

func (p *Plugin) Version() string {
	return xpb.Version()
}

func (p *Plugin) Description() string {
	return "The built-in pocketbase public_dir plugin included as an xpb plugin."
}

func (p *Plugin) Init(app core.App) error {

	if app, ok := app.(*pocketbase.PocketBase); ok {

		app.RootCmd.PersistentFlags().StringVar(
			&p.Dir,
			"publicDir",
			p.Dir,
			"the directory to serve static files",
		)

		app.RootCmd.PersistentFlags().BoolVar(
			&p.IndexFallback,
			"indexFallback",
			p.IndexFallback,
			"fallback the request to index.html on missing static path (eg. when pretty urls are used with SPA)",
		)

		app.RootCmd.ParseFlags(os.Args[1:])
	}

	// https://github.com/pocketbase/pocketbase/blob/f5c6b9652ffae8b2cffaa4d08430a2d955f2ad4e/examples/base/main.go#L106C1-L107C68
	//
	// static route to serves files from the provided public dir
	// (if publicDir exists and the route path is not already defined)
	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {
			if !e.Router.HasRoute(http.MethodGet, "/{path...}") {
				e.Router.GET("/{path...}", apis.Static(os.DirFS(p.Dir), p.IndexFallback))
			}
			return e.Next()
		},
		Priority: 999, // execute as latest as possible to allow users to provide their own route
	})

	return nil
}

// https://github.com/pocketbase/pocketbase/blob/f5c6b9652ffae8b2cffaa4d08430a2d955f2ad4e/examples/base/main.go#L124
//
// the default pb_public dir location is relative to the executable
func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "../pb_public")
}

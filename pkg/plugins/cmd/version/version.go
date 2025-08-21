package version

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbuilds/xpb"
	"github.com/pocketbuilds/xpb/pkg/plugins/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Register(&Plugin{})
}

var _ cmd.Command = (*Plugin)(nil)

type Plugin struct {
	ParentCmd *cobra.Command
}

func (p *Plugin) Name() string {
	return xpb.Version()
}

func (p *Plugin) Version() string {
	return xpb.Version()
}

func (p *Plugin) Description() string {
	return "A built-in xpb command for getting the current xpb version."
}

func (p *Plugin) SetParent(cmd *cobra.Command) {
	p.ParentCmd = cmd
}

func (p *Plugin) Init(app core.App) error {
	var cmd cobra.Command

	cmd.Use = "version"
	cmd.Short = "print xpb version number"
	cmd.Args = cobra.ExactArgs(0)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("xpb version", xpb.Version())
	}

	if p.ParentCmd == nil {
		p.ParentCmd = app.(*pocketbase.PocketBase).RootCmd
	}

	p.ParentCmd.AddCommand(&cmd)

	return nil
}

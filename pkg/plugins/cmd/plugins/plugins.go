package plugins

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbuilds/xpb"
	"github.com/pocketbuilds/xpb/pkg/plugins/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Register(&Plugin{})
}

type Plugin struct {
	ParentCmd *cobra.Command
}

func (p *Plugin) Name() string {
	return "plugins"
}

func (p *Plugin) Version() string {
	return xpb.Version()
}

func (p *Plugin) Description() string {
	return "A built-in xpb command for getting installed plugin information."
}

func (p *Plugin) SetParent(cmd *cobra.Command) {
	p.ParentCmd = cmd
}

func (p *Plugin) Init(app core.App) error {
	var cmd cobra.Command

	cmd.Use = "plugins [name]"
	cmd.Short = "Prints a list of plugins or prints the description the one specified"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		plugins := xpb.GetPlugins()

		if len(args) > 0 {
			arg := strings.ToLower(args[0])
			match := false
			for _, plugin := range plugins {
				name := plugin.Name()
				if len(arg) <= len(name) && arg == strings.ToLower(name)[:len(arg)] {
					if match {
						fmt.Println()
					} // Add extra newline if multiple plugins are being printed
					fmt.Printf("%s (%s)\n%s\n", plugin.Name(), plugin.Version(), plugin.Description())
					match = true
				}
			}
			if !match {
				fmt.Println("no plugins match this name")
			}
			return nil
		}

		fmt.Println("> Plugins")
		for _, plugin := range plugins {
			if version := plugin.Version(); version != "" {
				fmt.Printf("  - %s (%s)\n", plugin.Name(), version)
			} else {
				fmt.Printf("  - %s\n", plugin.Name())
			}
		}
		return nil
	}

	p.ParentCmd.AddCommand(&cmd)

	return nil
}

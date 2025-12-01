package plugins

import (
	"encoding/json"
	"fmt"
	"slices"
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
	var (
		cmd cobra.Command

		// flags
		jsonOutput bool
	)

	cmd.Use = "plugins [name]"
	cmd.Short = "Prints a list of plugins or prints the description the one specified"

	cmd.Flags().BoolVar(&jsonOutput, "json", jsonOutput, "format output as json")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		plugins := func() []xpb.Plugin {
			if len(args) > 0 {
				arg := strings.ToLower(args[0])
				return slices.Collect(func(yield func(xpb.Plugin) bool) {
					for _, p := range xpb.GetPlugins() {
						name := p.Name()
						if len(arg) <= len(name) && arg == strings.ToLower(name)[:len(arg)] {
							if !yield(p) {
								return
							}
						}
					}
				})
			}
			return xpb.GetPlugins()
		}()

		if len(plugins) == 0 && len(args) > 0 {
			return nil
		}

		if jsonOutput {
			pluginData := []map[string]string{}
			for _, plugin := range plugins {
				pluginData = append(pluginData, map[string]string{
					"name":        plugin.Name(),
					"version":     plugin.Version(),
					"description": plugin.Description(),
				})
			}
			jsonBytes, err := json.Marshal(pluginData)
			if err != nil {
				return err
			}
			fmt.Printf("%s", jsonBytes)
		} else {
			if len(plugins) == 0 {
				if len(args) > 0 {
					fmt.Println("no plugins match this name")
				}
				return nil
			}

			if len(plugins) == 1 {
				plugin := plugins[0]
				fmt.Printf("%s (%s)\n  %s\n", plugin.Name(), plugin.Version(), plugin.Description())
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
		}

		return nil
	}

	p.ParentCmd.AddCommand(&cmd)

	return nil
}

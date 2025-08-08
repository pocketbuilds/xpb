package plugin

import (
	"os"

	"github.com/PocketBuilds/xpb/pkg/templates"
	"github.com/spf13/cobra"
)

func init() {
	PluginCmd.AddCommand(initCmd)
}

var initCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [plugin name]",
		Short: "create new plugin project in current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			return templates.GeneratePluginDir(wd, templates.PluginTemplateData{
				Name: args[0],
			})
		},
	}
	return cmd
}()

package cmd

import "github.com/pocketbuilds/xpb/cmd/xpb/cmd/plugin"

func init() {
	rootCmd.AddCommand(plugin.PluginCmd)
}

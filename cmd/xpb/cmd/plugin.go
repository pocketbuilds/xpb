package cmd

import "github.com/PocketBuilds/xpb/cmd/xpb/cmd/plugin"

func init() {
	rootCmd.AddCommand(plugin.PluginCmd)
}

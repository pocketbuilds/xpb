package cmd

import (
	"github.com/PocketBuilds/xpb"
	"github.com/spf13/cobra"
)

var rootCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "xpb",
		Short:   "xpb is a tool for building customized pocketbase executables",
		Version: xpb.Version(),
	}
	return cmd
}()

func Execute() error {
	return rootCmd.Execute()
}

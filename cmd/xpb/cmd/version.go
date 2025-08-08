package cmd

import (
	"fmt"

	"github.com/PocketBuilds/xpb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print xpb version number",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("xpb version", xpb.Version())
		},
	}
	return cmd
}()

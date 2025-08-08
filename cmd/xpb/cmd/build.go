package cmd

import "github.com/PocketBuilds/xpb/cmd/xpb/cmd/build"

func init() {
	rootCmd.AddCommand(build.BuildCmd)
}

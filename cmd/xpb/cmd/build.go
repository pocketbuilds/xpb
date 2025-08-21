package cmd

import "github.com/pocketbuilds/xpb/cmd/xpb/cmd/build"

func init() {
	rootCmd.AddCommand(build.BuildCmd)
}

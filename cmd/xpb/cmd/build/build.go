package build

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pocketbuilds/xpb/pkg/builder"
	"github.com/pocketbuilds/xpb/pkg/module"
	"github.com/spf13/cobra"
)

var BuildCmd = func() *cobra.Command {
	var (
		output  = "pocketbase"
		with    = []string{}
		tags    = []string{"defaults"}
		config  = ""
		arch    = runtime.GOARCH
		osArg   = runtime.GOOS
		ldflags = []string{}
		dir     = ""
	)

	cmd := &cobra.Command{
		Use:   "build <version>",
		Short: "Build a custom pocketbase",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := exec.LookPath("go")
			if err != nil {
				return fmt.Errorf("go toolchain not found: %w", err)
			}

			var opts []builder.BuilderOption

			opts = append(opts,
				builder.WithArch(arch),
				builder.WithOS(osArg),
			)

			if len(args) == 1 {
				opts = append(opts, builder.WithPbVersion(args[0]))
			} else if v := os.Getenv("XPB__PB_VERSION"); v != "" {
				opts = append(opts, builder.WithPbVersion(v))
			}

			if dir != "" {
				opts = append(opts, builder.WithBuildDir(dir))
			}

			if config != "" {
				f, err := os.Open(config)
				if err != nil {
					return err
				}
				defer f.Close()
				opts = append(opts, builder.FromJsonReader(f))
			}

			for _, pluginArg := range with {
				opts = append(opts, builder.WithNewPlugin(module.FromCliArg(pluginArg)))
			}

			opts = append(opts, builder.WithTags(tags...))
			opts = append(opts, builder.WithLdflags(ldflags...))

			b, err := builder.NewBuilder(opts...)
			if err != nil {
				return err
			}

			return b.BuildToFile(output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", output, "output filepath")
	cmd.Flags().StringArrayVar(&with, "with", with, "include plugin  (format: module[@version][=replacement])")
	cmd.Flags().StringArrayVar(&tags, "tag", tags, "build tags")
	cmd.Flags().StringVarP(&config, "config", "c", config, "path to config file")
	cmd.Flags().StringVar(&arch, "arch", arch, "build target architecture")
	cmd.Flags().StringVar(&osArg, "os", osArg, "build target operating system")
	cmd.Flags().StringArrayVar(&ldflags, "ldflag", ldflags, "ldflags")
	cmd.Flags().StringVar(&dir, "dir", dir, "the directory the builder will use to build the go project (for debugging mostly)")

	return cmd
}()

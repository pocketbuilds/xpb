package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/PocketBuilds/xpb"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func init() {
	xpb.Register(cmd)
}

func Register(plugin xpb.Plugin) {
	cmd.children = append(cmd.children, plugin)
}

var cmd = &Plugin{}

type Plugin struct {
	children []xpb.Plugin
}

func (p *Plugin) Name() string {
	return "cmd"
}

func (p *Plugin) Version() string {
	return xpb.Version()
}

func (p *Plugin) Description() string {
	desc := "The built-in xpb commands including:\n"
	for _, c := range p.children {
		desc += fmt.Sprintf("  %s.%s (%s)\n%s\n", p.Name(), c.Name(), c.Version(), c.Description())
	}
	return desc
}

func (p *Plugin) Init(app core.App) error {
	pb, ok := app.(*pocketbase.PocketBase)
	if !ok {
		return nil
	}

	var xpbCmd = &cobra.Command{
		Use:   "xpb",
		Short: "xpb is a tool for building customized pocketbase executables",
	}

	for _, sm := range p.children {
		if cmd, ok := sm.(any).(Command); ok {
			cmd.SetParent(xpbCmd)
		}
		sm.Init(app)
	}

	pb.RootCmd.AddCommand(xpbCmd)

	return nil
}

type Command interface {
	SetParent(cmd *cobra.Command)
}

func (p *Plugin) UnmarshalJSON(data []byte) (err error) {
	// get raw json configs
	var configs map[string]json.RawMessage
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return err
	}

	// unmarshal raw json configs into children plugins
	for _, child := range p.children {
		if conf, ok := configs[child.Name()]; ok {
			err = json.Unmarshal(conf, child)
			if err != nil {
				return err
			}
		}
	}

	// unmarshal base module config
	type alias Plugin
	return json.Unmarshal(data, (*alias)(p))
}

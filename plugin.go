package xpb

import (
	"errors"
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
)

type Plugin interface {
	Init(app core.App) error
	Name() string
	Version() string
	Description() string
}

var plugins = []Plugin{}
var names = map[string]struct{}{}

func Register(plugin Plugin) {
	name := plugin.Name()
	if _, exists := names[name]; exists {
		log.Fatalf("fatal error: duplicate plugin name \"%s\"", name)
	}
	names[name] = struct{}{}
	plugins = append(plugins, plugin)
}

func GetPlugins() []Plugin {
	return plugins
}

func InitPlugins(app core.App) error {
	var errs []error

	for _, plugin := range plugins {
		var err error
		func() {
			defer func() {
				switch r := recover().(type) {
				case error:
					err = r
				case nil:
					return
				default:
					err = fmt.Errorf("%v", r)
				}
			}()
			err = plugin.Init(app)
		}()
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"error from plugin %s: %w",
				plugin.Name(),
				err,
			))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("error(s) during load: %w", errors.Join(errs...))
	}

	return nil
}

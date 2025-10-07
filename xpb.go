package xpb

import (
	"github.com/pocketbase/pocketbase/core"
)

func Setup(app core.App) error {
	stages := []func(core.App) error{
		PreValidatePlugins,
		LoadConfig,
		ValidatePlugins,
		InitPlugins,
	}
	for _, runStage := range stages {
		if err := runStage(app); err != nil {
			return err
		}
	}
	return nil
}

package xpb

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func LoadConfig(app core.App) (err error) {

	var dir = filepath.Dir(os.Args[0])

	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	var defaultFilePath = filepath.Join(dir, "pocketbuilds.toml")

	var fp = defaultFilePath

	if app, ok := app.(*pocketbase.PocketBase); ok {
		rootCmd := app.RootCmd
		rootCmd.PersistentFlags().StringVar(
			&fp,
			"config",
			defaultFilePath,
			"path to pocketbuilds toml config file",
		)
		rootCmd.ParseFlags(os.Args[1:])
	}

	f, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && fp == defaultFilePath {
			f, err = os.Create(fp)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()

	_, err = toml.NewDecoder(f).Decode(&config{})
	if err != nil {
		return err
	}

	for _, p := range plugins {
		err = env.ParseWithOptions(p, env.Options{
			Prefix: "XPB__" + strings.ToUpper(p.Name()) + "__",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type config struct{}

func (c *config) UnmarshalTOML(data any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.UnmarshalJSON(jsonBytes)
}

func (c *config) UnmarshalJSON(data []byte) (err error) {
	// get raw json configs
	var configs map[string]json.RawMessage
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return err
	}

	// unmarshal raw json configs into plugins
	for _, p := range plugins {
		if conf, ok := configs[p.Name()]; ok {
			err = json.Unmarshal(conf, p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

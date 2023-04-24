package main

import (
	"embed"
	"os"

	"github.com/spf13/viper"
)

//go:embed configs/*.yaml
var configs embed.FS

func loadConfig(paths ...string) error {
	viper.SetConfigType("yaml")

	// load default config
	f, err := configs.Open("configs/default.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := viper.ReadConfig(f); err != nil {
		return err
	}

	// load local configs
	for _, p := range paths {
		if p == "" {
			continue
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := viper.MergeConfig(f); err != nil {
			return err
		}
	}

	return nil
}

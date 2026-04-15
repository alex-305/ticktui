package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/apple/pkl-go/pkl"
	"github.com/pkg/errors"
)

type Config struct {
	DefaultProjectID    string              `pkl:"defaultProjectID"`
	DefaultProjectColor string              `pkl:"defaultProjectColor"`
	Keybindings         map[string][]string `pkl:"keybindings"`
}

var configPath = filepath.Join(xdg.ConfigHome, "ticktui", "config.pkl")

func InitConfig() error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return errors.Wrap(err, "creating config directory")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
			return errors.Wrap(err, "writing default pkl config")
		}
	}

	return nil
}

func Load(ctx context.Context) (*Config, error) {
	if err := InitConfig(); err != nil {
		return nil, err
	}

	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating pkl evaluator")
	}
	defer evaluator.Close()

	var cfg Config
	if err := evaluator.EvaluateModule(ctx, pkl.FileSource(configPath), &cfg); err != nil {
		return nil, errors.Wrap(err, "evaluating pkl config")
	}

	return &cfg, nil
}

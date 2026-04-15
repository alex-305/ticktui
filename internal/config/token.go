package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
)

var tokenPath = filepath.Join(xdg.DataHome, "ticktui", "token")

func LoadToken() (string, error) {
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func SaveToken(token string) error {
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		return errors.Wrap(err, "creating token directory")
	}

	return os.WriteFile(tokenPath, []byte(token), 0600)
}

func DeleteToken() error {
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

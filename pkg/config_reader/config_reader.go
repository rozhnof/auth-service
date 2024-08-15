package config_reader

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYaml(configPath string, out any) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

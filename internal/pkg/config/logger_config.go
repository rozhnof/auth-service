package config

type Logger struct {
	Level string `yaml:"level" env-required:"true"`
	Path  string `yaml:"path"`
}

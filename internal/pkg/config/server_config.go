package config

type Server struct {
	Address string `yaml:"address" env-required:"true"`
}

package config

type Tracing struct {
	Name     string `yaml:"name"`
	Output   string `yaml:"output" env-required:"true"`
	Endpoint string `yaml:"endpoint"`
}

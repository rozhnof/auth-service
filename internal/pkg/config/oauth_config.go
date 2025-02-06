package config

type Google struct {
	RedirectURL string   `yaml:"redirect" env-required:"true"`
	Scopes      []string `yaml:"scopes"   env-required:"true"`
}

type OAuth struct {
	Google Google `yaml:"google" env-required:"true"`
}

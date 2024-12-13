package config

type Redis struct {
	Address  string `env:"REDIS_ADDRESS"  env-required:"true"`
	Port     int    `env:"REDIS_PORT"     env-required:"true"`
	User     string `env:"REDIS_USER"     env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB"       env-required:"true"`
}

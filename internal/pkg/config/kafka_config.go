package config

type Kafka struct {
	BrokerList []string `yaml:"brokers" env-required:"true"`
}

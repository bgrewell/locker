package config

type Configuration struct {
	LogLevel string `yaml:"log_level"`
	FailOpen bool   `yaml:"fail_open"`
}

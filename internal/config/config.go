package config

func ReadConfiguration() (*Configuration, error) {
	return &Configuration{
		LockFileLocation: "/var/lock/locker.lock",
		LogLevel:         "info",
		FailOpen:         false,
		GRPCAddress:      "localhost:5128",
		BypassUsers:      []string{"root"},
		BypassGroups:     []string{"admin"},
	}, nil
}

type Configuration struct {
	LockFileLocation string   `yaml:"lock_file_location"`
	LogLevel         string   `yaml:"log_level"`
	FailOpen         bool     `yaml:"fail_open"`
	GRPCAddress      string   `yaml:"grpc_address"`
	BypassUsers      []string `yaml:"bypass_users"`
	BypassGroups     []string `yaml:"bypass_groups"`
}

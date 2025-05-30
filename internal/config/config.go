package config

type Config struct {
	CPUUsage    bool `mapstructure:"cpu_usage"`
	LoadAverage bool `mapstructure:"load_average"`
	DiskUsage   bool `mapstructure:"disk_usage"`
}

func NewConfig(configFile string) (*Config, error) {
	instance := &Config{}
	if err := LoadConfig(instance, configFile); err != nil {
		return nil, err
	}

	return instance, nil
}

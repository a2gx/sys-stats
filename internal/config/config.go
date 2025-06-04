package config

type Config struct {
	CPUUsage    bool `mapstructure:"cpu_usage"`
	LoadAverage bool `mapstructure:"load_average"`
	DiskUsage   bool `mapstructure:"disk_usage"`

	GRPC struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"grpc"`
}

func NewConfig(configFile string) (*Config, error) {
	instance := &Config{}
	if err := LoadConfig(instance, configFile); err != nil {
		return nil, err
	}

	//instance.GRPC.Host = "0.0.0.0"
	//instance.GRPC.Port = 50051

	return instance, nil
}

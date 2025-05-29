package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type TestConfigNested struct {
	Server   TestConfigNestedData `mapstructure:"server"`
	Database TestConfigNestedData `mapstructure:"database"`
}
type TestConfigNestedData struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

const nestedFileContent = `
server:
 host: "server.host"
 port: "8081"
database:
 host: "database.host"
 port: "5432"
`

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		env      map[string]string
		expected TestConfig
		wantErr  bool
	}{
		{
			name:   "loading file",
			config: "host: one.text\nport: 1234\ndatabase: test_1",
			expected: TestConfig{
				Host:     "one.text",
				Port:     1234,
				Database: "test_1",
			},
		},
		{
			name: "loading envelopment",
			env: map[string]string{
				"HOST":     "two.test",
				"PORT":     "4321",
				"DATABASE": "test_2",
			},
			expected: TestConfig{
				Host:     "two.test",
				Port:     4321,
				Database: "test_2",
			},
		},
		{
			name:   "combined config data",
			config: "host: three.text\nport: 1234\ndatabase: test",
			env: map[string]string{
				"PORT":     "5432",
				"DATABASE": "db_test",
			},
			expected: TestConfig{
				Host:     "three.text",
				Port:     5432,
				Database: "db_test",
			},
		},
		{
			name:     "invalid file",
			config:   "host: localhost\nport: non-a-number\ndatabase: test",
			expected: TestConfig{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset flags before test
			pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

			// create config file
			var configFile string
			if tt.config != "" {
				tmpFile, err := os.CreateTemp("", "config-*.yaml")
				require.NoError(t, err)
				defer os.Remove(tmpFile.Name())

				_, err = tmpFile.Write([]byte(tt.config))
				require.NoError(t, err)
				configFile = tmpFile.Name()
			}

			// environment variables
			savedEnv := make(map[string]string)
			for k, v := range tt.env {
				if old, exists := os.LookupEnv(k); exists {
					savedEnv[k] = old
				}
				os.Setenv(k, v)
			}

			// cleanup after test
			defer func() {
				for k := range tt.env {
					if old, exists := savedEnv[k]; exists {
						os.Setenv(k, old) // restore previous value
					} else {
						os.Unsetenv(k)
					}
				}
			}()

			// load configuration
			var conf TestConfig
			err := LoadConfig(&conf, configFile)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to unmarshal")
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, conf)
		})
	}
}

func TestLoadConfigNestedStruct(t *testing.T) {
	// reset flags before test
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// create config file
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(nestedFileContent))
	require.NoError(t, err)

	os.Setenv("SERVER_HOST", "server.env")
	os.Setenv("SERVER_PORT", "9090") // Change port value via ENV
	os.Setenv("DATABASE_HOST", "database.env")
	os.Setenv("DATABASE_PORT", "6000") // Similarly for database

	var conf TestConfigNested
	err = LoadConfig[TestConfigNested](&conf, tmpFile.Name())

	expected := TestConfigNested{
		Server: TestConfigNestedData{
			Host: "server.env", // Should come from ENV
			Port: 9090,         // Should come from ENV
		},
		Database: TestConfigNestedData{
			Host: "database.env",
			Port: 6000, // Should come from ENV
		},
	}

	require.NoError(t, err)
	require.Equal(t, expected, conf)
}

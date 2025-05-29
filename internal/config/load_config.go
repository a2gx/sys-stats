package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setDefaultEnv(prefix string, cfg interface{}) {
	valueOf, typeOf := reflect.ValueOf(cfg), reflect.TypeOf(cfg)

	if valueOf.Kind() == reflect.Ptr {
		valueOf, typeOf = valueOf.Elem(), typeOf.Elem()
	}

	for i := 0; i < typeOf.NumField(); i++ {
		fieldValue, fieldType := valueOf.Field(i), typeOf.Field(i)
		key := fieldType.Tag.Get("mapstructure")

		if prefix != "" {
			key = prefix + "." + key
		}

		if fieldValue.Kind() == reflect.Struct {
			setDefaultEnv(key, fieldValue.Addr().Interface())
		} else {
			viper.SetDefault(key, fieldValue.Interface())
		}
	}
}

func LoadConfig[T any](cfg *T, filepath string) error {
	viper.Reset()

	viper.SetConfigFile(filepath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set default environment variables
	setDefaultEnv("", cfg)

	// try to read config from file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Printf("Cannot read the configuration from the file: %s", err)
		}
	}

	// bind command line flags to configuration
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return fmt.Errorf("failed to bind command line flags: %w", err)
	}

	// unmarshal configuration into struct
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// config struct holds various configuration options.
type config struct {
	ApiEndpoint    string `yaml:"api"`
	LoginEndpoint  string `yaml:"login"`
	VerifyEndpoint string `yaml:"verify"`
}

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

// LoadConfig returns a configured viper instance
func LoadConfig() *config {
	v := readViperConfig("GO-HASTILY")
	conf := &config{}

	if err := v.Unmarshal(conf); err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// global defaults
	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	// read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return v
}

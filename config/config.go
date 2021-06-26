package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Config struct {
	Port   string `yaml:"port"`
	Log    *Log   `yaml:"log"`
	APIKey string `yaml:"apikey"`
}

const (
	portConfig      = "port"
	logLevelConfig  = "log.level"
	logFormatConfig = "log.format"
	apiKeyConfig    = "apikey"

	defaultPort      = "8080"
	defaultLogLevel  = "INFO"
	defaultLogFormat = ""
	defaultAPIKey    = ""
)

var (
	port      = pflag.StringP("port", "p", defaultPort, "The port the api will be served on")
	logLevel  = pflag.StringP("log-level", "", defaultLogLevel, "Sets the log level for the application (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)")
	logFormat = pflag.StringP("log-format", "", defaultLogFormat, "Sets the log output format")
	apiKey    = pflag.StringP("apikey", "", defaultLogFormat, "Your api key for openweathermap.org")
)

func NewConfig(cfgFile string) *Config {
	// Set the configuration file details
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/config/") // path to look for the config file in (docker container path)
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		fmt.Printf("Loading config from: %s\n", viper.ConfigFileUsed())
	}

	// Set default values for viper and bind the configuration together
	viper.SetDefault(portConfig, defaultPort)

	viper.BindEnv(logLevelConfig, "LOG_LEVEL")
	viper.SetDefault(logLevelConfig, defaultLogLevel)

	viper.BindEnv(logFormatConfig, "LOG_FORMAT")
	viper.SetDefault(logFormatConfig, defaultLogFormat)

	viper.SetDefault(apiKeyConfig, defaultAPIKey)

	// Process the configuration
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error as we will also pull from flags and env variables
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}
	viper.AutomaticEnv()

	overrideValuesFromFlags()

	// Pretty print viper settings
	fmt.Println("----------------------")
	fmt.Println("Parsed Config Settings")
	fmt.Println("----------------------")
	settings, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	fmt.Println(string(settings))
	fmt.Println("----------------------")

	// Build and return the configuration object
	config := &Config{
		Port: viper.GetString(portConfig),
		Log: &Log{
			Level:  viper.GetString(logLevelConfig),
			Format: viper.GetString(logFormatConfig),
		},
		APIKey: viper.GetString(apiKeyConfig),
	}

	return config
}

func overrideValuesFromFlags() {
	pflag.VisitAll(func(f *pflag.Flag) {
		// When a flag is provided on the commandline by a user, we will override the
		// discovered viper settings or any defaulted values.
		if f.Changed {
			// fmt.Printf("%s has been changed to %s\n", f.Name, f.Value)
			viper.Set(strings.ReplaceAll(f.Name, "-", "."), f.Value)
		}
	})
}

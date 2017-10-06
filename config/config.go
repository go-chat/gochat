package config

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Log contains logging settings.
type Log struct {
	Level string `yaml:"level"`
	// TODO: add external services (airbake, papertrail, slack, sentry...)
}

// Setup setups logrus log level, formatters... with the given settings in Log instance.
func (l *Log) Setup() (err error) {
	level, err := logrus.ParseLevel(l.Level)
	if err != nil {
		return
	}
	logrus.SetLevel(level)
	return
}

// SQL stores information like driver and datasource
type SQL struct {
	DriverName string
	DataSource string
}

// Config holds all settings of bank server.
type Config struct {
	Log *Log `yaml:"log"`
	SQL *SQL `yaml:"sql"`
}

// GetDefault creates new Config instance with all default settings.
func GetDefault() *Config {
	level := "info"

	if os.Getenv("GOCHAT_DEVELOPMENT") == "1" {
		level = "debug"
	}

	return &Config{
		Log: &Log{
			Level: level,
		},
		SQL: &SQL{
			DriverName: "postgres",
			DataSource: "postgresql://roach1@localhost:26257/gochat?sslmode=disable&connect_timeout=10",
		},
	}
}

// readFromFile reads the settings from the config file. If the config
// file isn't supplied, it uses all default values.
func readFromFile(filePath string) (cfg *Config, err error) {

	lf := logrus.Fields{
		"func":     "config.readFromFile",
		"filePath": filePath,
	}

	logrus.WithFields(lf).Debug("Load config")

	configBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.WithFields(lf).Error("Failed to load config file")
		return
	}
	cfg = GetDefault()
	err = yaml.Unmarshal(configBytes, cfg)
	return
}

// New reads the configs from arguments, validates and returns the
// config.
func New(cfgFile string) (cfg *Config, err error) {

	cfg, err = readFromFile(cfgFile)
	if err != nil {
		return
	}
	err = cfg.Log.Setup()
	return
}

// SetDevelopment set os GOCHAT_DEVELOPMENT to detect and enable/disable debug level
func SetDevelopment() error {
	return os.Setenv("GOCHAT_DEVELOPMENT", "1")
}

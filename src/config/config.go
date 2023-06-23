package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	App ApplicationConfig
}

func NewConfig() *Config {
	config := &Config{}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetDefault("config-path", "./config.yaml")
	configPath := viper.GetString("config-path")
	if _, err := os.Stat(configPath); err == nil {
		viper.SetConfigType(filepath.Ext(configPath)[1:])
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("failed to read config (%s): %+v", configPath, err)
		}
	}

	config.App.parse()

	return config
}

type ApplicationConfig struct {
	ServiceName    string
	WebPort        uint64
	HttpPathPrefix string

	OutgoingRequestTimeout     time.Duration
	NetworkBlockTime           time.Duration
	GetBlocksBatchSize         int
	OrphanPreventionBlockCount int

	DefaultStartingBlockNum  int
	DefaultStartingBlockHash string

	LogLevel string
}

func (c *ApplicationConfig) parse() {
	c.ServiceName = viperGetOrDefault("app.service-name", "web3-api-proxy")
	c.WebPort = viperGetOrDefaultUint64("app.web-port", 8000)
	c.HttpPathPrefix = viperGetOrDefault("app.http-path-prefix", "")
	c.OutgoingRequestTimeout = viperGetOrDefaultTimeDuration("app.outgoing-request-timeout", "15s")
	c.NetworkBlockTime = viperGetOrDefaultTimeDuration("app.network-block-time", "12s")
	c.GetBlocksBatchSize = viperGetOrDefaultInt("app.get-blocks-batch-size", 5)
	c.OrphanPreventionBlockCount = viperGetOrDefaultInt("app.orphan-prevention-block-count", 10)
	c.DefaultStartingBlockNum = viperGetOrDefaultInt("app.default-starting-block-num", 17539747)
	c.DefaultStartingBlockHash = viperGetOrDefault("app.default-starting-block-hash", "0x206b9c2a16f0774b64b2b682683db378aacdd619b532f391633181b99abb1a41")
	c.LogLevel = viperGetOrDefault("app.log-level", "debug")
}

func viperGetOrDefault(key string, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func viperGetOrDefaultInt(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}

func viperGetOrDefaultUint64(key string, defaultValue uint64) uint64 {
	viper.SetDefault(key, defaultValue)
	return viper.GetUint64(key)
}

func viperGetOrDefaultTimeDuration(key string, defaultValue string) time.Duration {
	viper.SetDefault(key, defaultValue)
	d, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		log.Fatalf("provided value '%s' cannot be transformed to [time.Duration]", viper.GetString(key))
	}
	return d
}

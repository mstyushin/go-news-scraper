package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	VersionDate = "0001-01-01 00:00:00"
	Version     = "dev"
	Hash        = "COMMIT ID"
)

// DefaultConfig возвращает конфигурацию по умолчанию.
// будет использоваться в случае когда config.yaml не найден
func DefaultConfig() *Config {
	return &Config{
		HttpPort:      8081,
		RequestPeriod: 10 * time.Minute,
		DBConnString:  "postgres://postgres@localhost:5432/news?sslmode=disable",
		RSSFeeds: []string{
			"https://habr.com/ru/rss/hub/go/all/?fl=ru",
			"https://habr.com/ru/rss/best/daily/?fl=ru",
			"https://cprss.s3.amazonaws.com/golangweekly.com.xml",
		},
	}
}

// конфигурация приложения, подразумевается yaml-формат
type Config struct {
	HttpPort      int           `yaml:"http_port,omitempty"`
	RequestPeriod time.Duration `yaml:"request_period,omitempty"`
	DBConnString  string        `yaml:"db_conn_string,omitempty"`
	RSSFeeds      []string      `yaml:"rss_feeds,omitempty"`
}

func VersionString() string {
	return fmt.Sprintf("Version: %s Commit: %s BuildDate: %s", Version, Hash, VersionDate)
}

func New() (*Config, error) {
	var config *Config

	var configPath string
	var printConfig bool
	var printVersion bool
	flag.StringVar(&configPath, "config", "./config.yaml", "path to a YAML config file")
	flag.BoolVar(&printConfig, "print-config", false, "print loaded config")
	flag.BoolVar(&printVersion, "version", false, "print build version")
	flag.Parse()

	if printVersion {
		fmt.Println(VersionString())
		return nil, nil
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("not found config file at %s, using defaults\n", configPath)
		config = DefaultConfig()
	} else {
		log.Printf("reading config at %s\n", configPath)
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return nil, err
	}

	if printConfig {
		yamlData, _ := yaml.Marshal(&config)
		fmt.Println(string(yamlData))
		return nil, nil
	}

	return config, nil
}

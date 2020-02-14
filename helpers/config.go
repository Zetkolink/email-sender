package helpers

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig() Config {
	cfg := Config{}
	cfg.loadConfig()
	return cfg
}

type Config struct {
	Db        DbConfig     `yaml:"db"`
	Smtp      SmtpConfig   `yaml:"smtp"`
	Rb        RabbitConfig `yaml:"rb"`
	LogLevel  string       `yaml:"logLevel"`
	LogFormat string       `yaml:"logFormat"`
}

type DbConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Schema string `yaml:"schema"`
}

type SmtpConfig struct {
	Identity string `yaml:"identity"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
}

type RabbitConfig struct {
	Amqp    string `yaml:"amqp"`
	Channel string `yaml:"channel"`
}

func (c *Config) loadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
}

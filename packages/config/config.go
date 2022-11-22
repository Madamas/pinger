package config

import (
	"time"

	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type Client struct {
	Timeout         string `yaml:"timeout" required:"true"`
	TimeoutDuration time.Duration
}

type PingerTarget struct {
	Host  string
	Port  int
	Route string
}

type Pinger struct {
	Receiver         int64
	Interval         string
	IntervalDuration time.Duration
	Targets          []PingerTarget
}

type Bot struct {
	Token string `env:"BOT_TOKEN" required:"true"`
	Debug bool
}

type Config struct {
	Client Client
	Pinger Pinger
	Bot    Bot
}

func (c Client) parse() (Client, error) {
	var err error
	c.TimeoutDuration, err = time.ParseDuration(c.Timeout)

	if err != nil {
		panic(errors.Wrap(err, "Couldn't read client timeout duration"))
	}

	return c, err
}

func (c Config) parse() (Config, error) {
	var err error
	c.Client, err = c.Client.parse()

	return c, err
}

func Populate() (Config, error) {
	conf := Config{}

	err := configor.New(&configor.Config{
		ErrorOnUnmatchedKeys: true,
		Verbose:              true,
	}).Load(&conf, "config/config.yaml")

	if err != nil {
		return conf, err
	}

	return conf.parse()
}

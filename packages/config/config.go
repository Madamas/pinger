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

type Storage struct {
	Url                    string `required:"true"`
	Database               string `required:"true"`
	InfoCollection         string `yaml:"infoCollection" required:"true"`
	ConnectTimeout         string `yaml:"connectTimeout" required:"true"`
	WriteTimeout           string `yaml:"writeTimeout" required:"true"`
	ConnectTimeoutDuration time.Duration
	WriteTimeoutDuration   time.Duration
}

type Config struct {
	Client  Client
	Pinger  Pinger
	Bot     Bot
	Storage Storage
}

func (s Storage) parse() (Storage, error) {
	var err error
	s.ConnectTimeoutDuration, err = time.ParseDuration(s.ConnectTimeout)
	if err != nil {
		panic(errors.Wrap(err, "Couldn't read connect timeout duration"))
	}

	s.WriteTimeoutDuration, err = time.ParseDuration(s.WriteTimeout)
	if err != nil {
		panic(errors.Wrap(err, "Couldn't read write timeout duration"))
	}

	return s, err
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

	if err != nil {
		return c, err
	}

	c.Storage, err = c.Storage.parse()

	if err != nil {
		return c, err
	}

	c.Pinger, err = c.Pinger.parse()

	return c, err
}

func (p Pinger) parse() (Pinger, error) {
	var err error

	p.IntervalDuration, err = time.ParseDuration(p.Interval)

	if err != nil {
		panic(errors.Wrap(err, "Couldn't read pinger interval duration"))
	}

	return p, err
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

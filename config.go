package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Resolution       string   `yaml:"resolution"`
	Width            int      `yaml:"width"`
	Height           int      `yaml:"height"`
	Fullscreen       bool     `yaml:"fullscreen"`
	Debug            bool     `yaml:"debug"`
	BallSpeed        float64  `yaml:"ballSpeed"`
	BallAcceleration float64  `yaml:"ballAcceleration"`
	ListenInterfaces []string `yaml:"listenInterfaces"`

	path string
}

func (c *Config) Load() *Config {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return c
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Fatal(err)
	}

	res := strings.Split(c.Resolution, "x")
	if len(res) != 2 {
		log.Fatalf("invalid resolution: %s", c.Resolution)
	}

	c.Width, err = strconv.Atoi(res[0])
	if err != nil {
		log.Fatalf("invalid width: %s", res[0])
	}

	c.Height, err = strconv.Atoi(res[1])
	if err != nil {
		log.Fatalf("invalid height: %s", res[1])
	}

	return c
}

func NewConfig() *Config {
	return &Config{
		Resolution:       "1280x720",
		Width:            1280,
		Height:           720,
		Fullscreen:       false,
		BallSpeed:        5,
		BallAcceleration: 1.1,

		path: "config.yml",
	}
}

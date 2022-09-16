package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Resolution             string  `yaml:"resolution"`
	Width                  int     `yaml:"width"`
	Height                 int     `yaml:"height"`
	Fullscreen             bool    `yaml:"fullscreen"`
	Debug                  bool    `yaml:"debug"`
	BallSpeed              float64 `yaml:"ballSpeed"`
	BallAcceleration       float64 `yaml:"ballAcceleration"`
	ListenPort             int     `yaml:"listenPort"`
	Offline                bool    `yaml:"offline"`
	AnnounceFrequency      int     `yaml:"announceFrequency"`
	StateSyncFrequency     int     `yaml:"stateSyncFrequency"`
	StateFullSyncFrequency int     `yaml:"stateFullSyncFrequency"`
	DisableVSync           bool    `yaml:"disableVSync"`
	TicksPerSecond         int     `yaml:"ticksPerSecond"`

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
		Resolution:             "1280x720",
		Width:                  1280,
		Height:                 720,
		Fullscreen:             false,
		BallSpeed:              5,
		BallAcceleration:       1.1,
		Offline:                false,
		AnnounceFrequency:      3,
		StateSyncFrequency:     30,
		StateFullSyncFrequency: 2,
		DisableVSync:           false,
		TicksPerSecond:         60,

		path: "config.yml",
	}
}

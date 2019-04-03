package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type config struct {
	SourceFolder       string   `json:"source-folder"`
	ConfigFileSuffixes []string `json:"config-file-suffixes"`
	Rules              []string `json:"rules"`
}

var defaultConfig = config{
	SourceFolder:       "~/IdeaProjects/cms/play",
	ConfigFileSuffixes: []string{".yaml", ".properties"},
}

func mustParseConfig() config {
	config, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("config", config).Info("config: parsed config")
	return config
}

func parseConfig() (config, error) {
	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "", "Config file")
	flag.Parse()

	if configFile == "" {
		flag.Usage()
		return defaultConfig, fmt.Errorf("config file is required")
	}

	byts, err := ioutil.ReadFile(configFile)
	if err != nil {
		return defaultConfig, fmt.Errorf("couldn't read config file. err=%v", err)
	}

	err = json.Unmarshal(byts, &defaultConfig)
	if err != nil {
		return defaultConfig, fmt.Errorf("couldn't unmarshal config file. err=%v", err)
	}

	return defaultConfig, nil
}

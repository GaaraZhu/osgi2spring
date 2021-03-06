package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type config struct {
	SourceFolder      string   `json:"source-folder"`
	FilesToBeExcluded []string `json:"files-to-be-excluded"`
	Rules             []string `json:"rules"`
	Mode              string   `json:"mode"`

	filesToBeExcluded map[string]string
}

var defaultConfig = config{
	SourceFolder: "~/IdeaProjects/cms/play",
	FilesToBeExcluded: []string{"RedissonConfiguration.java", "PreBid.java", "EndpointSentryConfig.java", "RemoteJcrEventListener.java",
		"TagManagerConfig.java", "OsgiServiceReferenceResolver.java", "SegmentConfig.java", "PlayApplication.java", "JerseyConfig.java",
		"ConsumerAuthenticationFilter.java"},
	Mode:              "static",
	filesToBeExcluded: make(map[string]string),
}

func mustParseConfig() config {
	config, err := parseConfig()
	if err != nil {
		panic(err)
	}
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

	for _, f := range defaultConfig.FilesToBeExcluded {
		defaultConfig.filesToBeExcluded[f] = ""
	}

	return defaultConfig, nil
}

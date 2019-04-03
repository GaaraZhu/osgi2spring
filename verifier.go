package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := mustParseConfig()

	// 1. get all rules
	var rules []rule
	for _, r := range cfg.Rules {
		rules = append(rules, ruleRegistry[r])
	}

	log.WithField("rules", cfg.Rules).Info("scanner: files to be scanned:")

	// 2. get all source files
	var filesToBeScanned []string
	var filesNotMeetRules map[string][]string = make(map[string][]string)
	err := filepath.Walk(cfg.SourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{"error": err, "folder": cfg.SourceFolder}).Error("scanner: failed to walk through source folder")
			return err
		}

		if !info.IsDir() {
			filesToBeScanned = append(filesToBeScanned, path)
		}
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{"error": err, "folder": cfg.SourceFolder}).Error("scanner: quit as it failed to walk through source folder")
		return
	}

	log.WithField("file-count", len(filesToBeScanned)).Info("scanner: files to be scanned:")

	// 3. verify rule one by one
	for _, file := range filesToBeScanned {
		source, e := extractSource(file)
		if e != nil {
			log.WithField("file", file).Error("scanner: failed to extract source code")
			return
		}

		for _, rule := range rules {
			isMet, err := rule.isMet(source)
			if err != nil {
				log.WithFields(log.Fields{"file": file, "rule": rule}).Error("scanner: failed to verify rule")
				return
			}

			if !isMet {
				if _, ok := filesNotMeetRules[file]; ok {
					filesNotMeetRules[file] = append(filesNotMeetRules[file], rule.getName())
				} else {
					filesNotMeetRules[file] = []string{rule.getName()}
				}
			}
		}
	}

	// 4. verification summary
	log.WithField("filesNotMeetRules", filesNotMeetRules).Info("scanner: finished scanning. Files to be updated:")
}

func extractSource(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {

		return "", err
	}

	return string(data), nil
}

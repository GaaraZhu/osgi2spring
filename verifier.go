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

	log.WithField("rules", cfg.Rules).Info("verifier: rules to be verified:")

	// 2. get all source files
	var filesToBeScanned []string
	var filesNotMeetRules = make(map[string][]string)
	err := filepath.Walk(cfg.SourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{"error": err, "folder": cfg.SourceFolder}).Error("verifier: failed to walk through source folder")
			return err
		}

		if !info.IsDir() {
			filesToBeScanned = append(filesToBeScanned, path)
		}
		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{"error": err, "folder": cfg.SourceFolder}).Error("verifier: quit as it failed to walk through source folder")
		return
	}

	log.WithField("file-count", len(filesToBeScanned)).Info("verifier: total files to be scanned:")

	// 3. verify rule one by one
	for _, file := range filesToBeScanned {
		source, e := extractSource(file)
		if e != nil {
			log.WithField("file", file).Error("verifier: failed to extract source code")
			return
		}

		for _, rule := range rules {
			isMet, err := rule.isMet(source)
			if err != nil {
				log.WithFields(log.Fields{"file": file, "rule": rule}).Error("verifier: failed to verify rule")
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
	if len(filesNotMeetRules) > 0 {
		log.Info("verifier: verify result: failed -_-")
	} else {
		log.Info("verifier: verify result: success ^_^")
	}
	for k, v := range filesNotMeetRules {
		log.WithField("file", k).Info("verifier: file:")
		log.WithField("rule(s)", v).Info("verifier: broken rule(s):")
	}
}

func extractSource(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {

		return "", err
	}

	return string(data), nil
}

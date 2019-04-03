package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := mustParseConfig()
	var filesToBeScanned, filesToBeUpdated []string

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

	for _, file := range filesToBeScanned {
		needsUpdate, e := scanSingleFile(file)
		if e != nil {
			return
		}

		if needsUpdate {
			filesToBeUpdated = append(filesToBeUpdated, file)
		}
	}

	log.WithField("filesToBeUpdated", filesToBeUpdated).Info("scanner: finished scanning. Files to be updated:")
}

func scanSingleFile(path string) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithField("path", path).Error("scanner: failed to read file")
		return false, err
	}

	source := string(data)
	if strings.Contains(source, "@Component") && !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return true, nil
	}

	return false, nil
}

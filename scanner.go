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

	source := string(data)                            // get the source code string
	source = strings.Join(strings.Fields(source), "") // removes all spaces/taps/new lines

	if !strings.Contains(source, "@Component") {
		return false, nil // skip for non OSGIN components
	}

	//Rule No.1: spring component annotation should be present for OSGI components
	if !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return true, nil
	}

	//Rule No.2: no-args activate method should be invoked after the construction
	if strings.Contains(source, "voidactivate()") && !strings.Contains(source, "@PostConstruct") {
		return true, nil
	}

	//Rule No.3: deactive method should be invoked before the object is destroyed
	if strings.Contains(source, "@Deactivate") && !strings.Contains(source, "@PreDestroy") {
		return true, nil
	}

	//Rule No.4 OSGI injected properties should be injected by spring
	if strings.Contains(source, "@Reference") {
		sr := strings.Split(path, "\\")
		className := strings.Split(sr[len(sr)-1], ".java")[0]
		sr = strings.Split(source, "@Reference")
		for i, r := range sr {
			if i == 0 {
				continue
			}
			subsequentCode := strings.Split(r, ";")[0]
			if !(strings.Contains(subsequentCode, "@Autowired") ||
				strings.Contains(subsequentCode, "@Qualifier") ||
				strings.Contains(subsequentCode, "@Resource") ||
				strings.Contains(subsequentCode, "@Inject") ||
				strings.Contains(source, "@Autowiredpublic"+className)) {
				return true, nil
			}
		}

		return false, nil
	}

	return false, nil
}

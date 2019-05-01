package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	c "github.com/logrusorgru/aurora"
	"gopkg.in/cheggaaa/pb.v1"
)

var toolName = c.Yellow("osgi2spring").Bold()

func main() {
	cfg := mustParseConfig()
	fmt.Printf("%s: Java source folder:\n\t * %s\n", toolName, cfg.SourceFolder)

	// 1. get all rules
	rules := extractRules(cfg)
	fmt.Printf("%s: Rule(s) to be verified: \n", toolName)
	for _, rule := range rules {
		fmt.Printf("\t * %s: %s\n", rule.getName(), rule.getDescription())
	}
	fmt.Printf("\n")

	// 2. get all source files
	sourceFiles, err := extractSourceFiles(cfg)
	if err != nil {
		fmt.Printf("%s: Quit as it failed to walk through source folder due to: %v\n", toolName, err)
		return
	}
	fmt.Printf("%s: Total files to be verified: %v\n", toolName, c.Green(len(sourceFiles)))

	// 3. verify files&rules one by one
	filesNotMeetRules, err := process(sourceFiles, rules, cfg.Mode)
	if err != nil {
		fmt.Printf("%s: Quitting due to process error, %v", toolName, err)
		return
	}

	if len(filesNotMeetRules) == 0 {
		return
	}

	// 4. verification summary
	// sort map keys(file names)
	var fileNames []string
	for k := range filesNotMeetRules {
		fileNames = append(fileNames, k)
	}
	sort.Strings(fileNames)

	fmt.Printf("%s: Result details: \n", toolName)
	for _, k := range fileNames {
		fmt.Printf("%s: %s\n", c.Cyan("File"), k)
		fmt.Printf("%s: %s\n\n", c.Cyan("Broken rule(s)"), strings.Join(filesNotMeetRules[k], ", "))
	}
}

func process(files []string, rules []rule, mode string) (map[string][]string, error) {
	var (
		err               error
		beanPayload       string
		filesNotMeetRules = make(map[string][]string)
		bar               = pb.StartNew(len(files))
	)

	if mode != "static" {
		beanPayload, err = getBeansPayload()
		if err != nil {
			fmt.Printf("%s: Failed to get runtime beans due to: %v\n", toolName, err)
			return filesNotMeetRules, err
		}
	}
	for _, file := range files {
		source, err := extractSource(file)
		if err != nil {
			fmt.Printf("%s: Failed to extract source code due to: %v\n", toolName, err)
			return filesNotMeetRules, err
		}

		for _, rule := range rules {
			if mode == "static" {
				isMet, err := rule.isMetStaticly(source)
				if err != nil {
					fmt.Printf("%s: Failed to verify rule due to: %v\n", toolName, file)
					return filesNotMeetRules, err
				}

				if !isMet {
					if _, ok := filesNotMeetRules[file]; ok {
						filesNotMeetRules[file] = append(filesNotMeetRules[file], rule.getName())
					} else {
						filesNotMeetRules[file] = []string{rule.getName()}
					}
				}
				continue
			}

			isMet, err := rule.isMetRuntimely(source, beanPayload)
			if err != nil {
				fmt.Printf("%s: Failed to verify rule due to: %v\n", toolName, file)
				return filesNotMeetRules, err
			}

			if !isMet {
				if _, ok := filesNotMeetRules[file]; ok {
					filesNotMeetRules[file] = append(filesNotMeetRules[file], rule.getName())
				} else {
					filesNotMeetRules[file] = []string{rule.getName()}
				}
			}

		}
		bar.Increment()
	}
	var finishMessage = fmt.Sprintf("%s: Result summary: success ^_^\n\n", toolName)
	if len(filesNotMeetRules) > 0 {
		finishMessage = fmt.Sprintf("%s: Result summary: %v files failed -_-\n", toolName, c.Red(len(filesNotMeetRules)))
	}
	bar.FinishPrint(finishMessage)

	return filesNotMeetRules, nil
}

func extractRules(cfg config) []rule {
	var rules []rule
	for _, r := range cfg.Rules {
		rules = append(rules, ruleRegistry[r])
	}

	return rules
}

func extractSourceFiles(cfg config) ([]string, error) {
	var sourceFiles []string
	err := filepath.Walk(cfg.SourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%s: Failed to walk through source folder due to: %v", toolName, err)
		}

		if !info.IsDir() && !needToBeExclude(info.Name(), cfg) {
			sourceFiles = append(sourceFiles, path)
		}
		return nil
	})

	if err != nil {

		return []string{}, err
	}

	return sourceFiles, nil
}

func extractSource(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {

		return "", err
	}

	return string(data), nil
}

func needToBeExclude(fileName string, cfg config) bool {
	for _, f := range cfg.FilesToBeExcluded {
		if f == fileName {
			return true
		}
	}
	return false
}

func getBeansPayload() (string, error) {
	resp, err := http.Get("http://localhost:9202/actuator/beans")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

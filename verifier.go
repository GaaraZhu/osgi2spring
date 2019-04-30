package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	cfg := mustParseConfig()
	fmt.Printf("verifier: verification configs: %v \n", cfg)

	// 1. get all rules
	rules := extractRules(cfg)

	// 2. get all source files
	sourceFiles, err := extractSourceFiles(cfg)
	if err != nil {
		fmt.Printf("verifier: quit as it failed to walk through source folder due to: %v \n", err)
		return
	}
	fmt.Printf("verifier: total files to be verified: %v \n", len(sourceFiles))

	// 3. verify files&rules one by one
	filesNotMeetRules, err := process(sourceFiles, rules)
	if err != nil {
		fmt.Printf("verifier: quitting due to process error, %v", err)
		return
	}

	// sort map keys(file names)
	var fileNames []string
	for k := range filesNotMeetRules {
		fileNames = append(fileNames, k)
	}
	sort.Strings(fileNames)

	// 4. verification summary
	fmt.Println("-------------------------------------------")
	fmt.Println("result details:")
	for _, k := range fileNames {
		fmt.Printf("file path: [%s], broken rule(s): %v \n", k, filesNotMeetRules[k])
	}

	if len(filesNotMeetRules) > 0 {
		fmt.Println()
		fmt.Printf("verified summary: %v files failed -_-", len(filesNotMeetRules))
	} else {
		fmt.Println("verified summary: success ^_^")
	}
}

func process(files []string, rules []rule) (map[string][]string, error) {
	bar := pb.StartNew(len(files))

	var filesNotMeetRules = make(map[string][]string)
	for _, file := range files {
		source, err := extractSource(file)
		if err != nil {
			fmt.Printf("verifier: failed to extract source code due to: %v \n", err)
			return filesNotMeetRules, err
		}

		for _, rule := range rules {
			isMet, err := rule.isMet(source)
			if err != nil {
				fmt.Printf("verifier: failed to verify rule due to: %v \n", file)
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
	bar.FinishPrint("The End!")

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
			return fmt.Errorf("verifier: failed to walk through source folder due to: %v", err)
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

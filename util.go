package main

import (
	"strings"
)

// extracts all OSGI components from a signle Java source file
func extractOSGIComponents(source string) []string {
	if !strings.Contains(source, "@Component") {
		return []string{}
	}

	componentSplitResult := strings.Split(source, "@Component")
	var components []string
	for i := 1; i < len(componentSplitResult); i++ {
		component := extractSingleComponent(componentSplitResult[i])
		if component != "" {
			components = append(components, component)
		}
	}

	return components
}

func extractSingleComponent(sourceSplit string) string {
	tmp := strings.Split(sourceSplit, "{")
	sr := strings.Split(tmp[0], "public class ")
	if len(sr) < 2 {
		sr = strings.Split(tmp[0], "public static class ")
		if len(sr) < 2 {
			sr = strings.Split(tmp[0], "public abstract class ")
			if len(sr) < 2 {
				return "" // we are not supposed to annotate interfaces
			}
		}
	}

	component := strings.Split(sr[1], " ")[0]
	if strings.Contains(component, "<") {
		component = strings.Split(component, "<")[0]
	}

	return component
}

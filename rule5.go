package main

import "strings"

type rule5 struct{}

func (r rule5) getName() string {
	return "property-injected"
}

func (r rule5) getDescription() string {
	return "non-args constructor should be present for OSGI"
}

func (r rule5) isMet(source string) (bool, error) {
	if !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return true, nil // skip for non spring beans
	}

	sr := strings.Split(source, "public class ")
	if len(sr) < 2 {
		sr = strings.Split(source, "public abstract class ")
		if len(sr) < 2 {
			return true, nil // skip for interface
		}
	}

	className := strings.Split(sr[1], " ")[0]

	if (strings.Contains(source, "public "+className+"(@Value") || strings.Contains(source, "public "+className+" (@Value")) && (!(strings.Contains(source, "public "+className+"()") || strings.Contains(source, "public "+className+" ()"))) {
		return false, nil
	}

	return true, nil
}
package main

import "strings"

type rule1 struct{}

func (r rule1) getName() string {
	return "spring-containered"
}

func (r rule1) getDescription() string {
	return "OSGI components should be managed by spring container"
}

func (r rule1) isMet(source string) (bool, error) {
	if !strings.Contains(source, "@Component") {
		return true, nil // skip for non OSGIN components
	}

	if !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return false, nil
	}

	return true, nil
}

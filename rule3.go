package main

import "strings"

type rule3 struct{}

func (r rule3) getName() string {
	return "deactivation-invoked"
}

func (r rule3) getDescription() string {
	return "deactivation method should be invoked before object is destoryed"
}

func (r rule3) isMet(source string) (bool, error) {
	if !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return true, nil // skip for non spring components
	}

	if strings.Contains(source, "void deactivate") && !strings.Contains(source, "@PreDestroy") {
		return false, nil
	}

	return true, nil
}

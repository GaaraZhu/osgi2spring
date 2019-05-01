package main

import (
	"strings"
)

type rule3 struct{}

func (r rule3) getName() string {
	return "deactivation-invoked"
}

func (r rule3) getDescription() string {
	return "deactivation method should be invoked before object is destoryed"
}

func (r rule3) isMet(source string) (bool, error) {
	if !strings.Contains(source, "@Component") {
		return true, nil // skip for non OSGI components
	}

	if !strings.Contains(source, "@Deactivate") {
		return true, nil // skip if no deactivation
	}

	//non-args deactivation should be triggered automatically by using @PreDestroy
	if strings.Contains(source, "void deactivate()") {
		if !strings.Contains(source, "@PreDestroy") {
			return false, nil
		}
		return true, nil
	}

	//activation with args should be invoked manually after property injection
	if strings.Count(source, "deactivate(") < 2 {
		return false, nil
	}

	return true, nil
}

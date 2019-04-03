package main

import "strings"

type rule2 struct{}

func (r rule2) getName() string {
	return "activation-invoked"
}

func (r rule2) getDescription() string {
	return "activation method should be invoked after object construction"
}

func (r rule2) isMet(source string) (bool, error) {
	if !strings.Contains(source, "@org.springframework.stereotype.Component") {
		return true, nil // skip for non spring components
	}

	if strings.Contains(source, "@Activate") {
		//non-args activation should be triggered automatically by using @PostConstruct
		if strings.Contains(source, "void activate()") {
			if !strings.Contains(source, "@PostConstruct") {
				return false, nil
			}
			return true, nil
		}

		//args activation should take property injection into account and should be invoked manually
		source = strings.NewReplacer("deactivate", "").Replace(source)
		strings.Count(source, "activate(")
		if strings.Count(source, "activate(") < 2 {
			return false, nil
		}

	}

	return true, nil
}

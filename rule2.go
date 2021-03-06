package main

import "strings"

type rule2 struct{}

func (r rule2) getName() string {
	return "activation-invoked"
}

func (r rule2) getDescription() string {
	return "activation method should be invoked after object construction"
}

func (r rule2) isMetStaticly(source string) (bool, error) {
	if !strings.Contains(source, "@Component") {
		return true, nil // skip for non OSGI components
	}

	if !strings.Contains(source, "@Activate") {
		return true, nil // skip if no activation
	}

	//non-args activation should be triggered automatically by using @PostConstruct
	if strings.Contains(source, "void activate()") {
		if !strings.Contains(source, "@PostConstruct") {
			return false, nil
		}
		return true, nil
	}

	//activation with args should be invoked manually after property injection
	source = strings.NewReplacer("deactivate", "").Replace(source)
	if strings.Count(source, "activate(") < 2 {
		return false, nil
	}

	return true, nil
}

func (r rule2) isMetRuntimely(source, beanPayload string) (bool, error) {
	return true, nil
}

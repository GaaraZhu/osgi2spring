package main

type rule interface {
	getName() string

	getDescription() string

	isMetStaticly(source string) (bool, error)

	isMetRuntimely(source, beanPayload string) (bool, error)
}

var ruleRegistry = make(map[string]rule)

func init() {
	ruleRegistry["rule1"] = rule1{}
	ruleRegistry["rule2"] = rule2{}
	ruleRegistry["rule3"] = rule3{}
	ruleRegistry["rule4"] = rule4{}
	ruleRegistry["rule5"] = rule5{}
}

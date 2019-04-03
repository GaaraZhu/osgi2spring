package main

type rule interface {
	getName() string
	getDescription() string
	isMet(source string) (bool, error) //spaces, tabs and new lines should be removed from source before passing here
}

var ruleRegistry = make(map[string]rule)

func init() {
	ruleRegistry["rule1"] = rule1{}
	ruleRegistry["rule2"] = rule2{}
	ruleRegistry["rule3"] = rule3{}
	ruleRegistry["rule4"] = rule4{}
}

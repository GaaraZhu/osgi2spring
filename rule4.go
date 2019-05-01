package main

import "strings"

type rule4 struct{}

func (r rule4) getName() string {
	return "property-injected"
}

func (r rule4) getDescription() string {
	return "OSGI injected properties should be injected by spring"
}

func (r rule4) isMetStaticly(source string) (bool, error) {
	if !strings.Contains(source, "@Component") {
		return true, nil // skip for non OSGI beans
	}

	if strings.Contains(source, "@Reference") {
		sr := strings.Split(source, "public class ")
		if len(sr) < 2 {
			sr = strings.Split(source, "public abstract class ")
			if len(sr) < 2 {
				return true, nil // skip for interface
			}
		}
		className := strings.Split(sr[1], " ")[0]

		source = strings.Join(strings.Fields(source), "")
		sr = strings.Split(source, "@Reference")
		for i, r := range sr {
			if i == 0 {
				continue
			}
			subsequentCode := strings.Split(r, ";")[0]
			if !(strings.Contains(subsequentCode, "@Autowired") ||
				strings.Contains(subsequentCode, "@Resource") ||
				strings.Contains(subsequentCode, "@Inject") ||
				strings.Contains(source, "@Autowiredpublic"+className)) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (r rule4) isMetRuntimely(source, beanPayload string) (bool, error) {
	return true, nil
}

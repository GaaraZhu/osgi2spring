package main

import (
	"fmt"
	"testing"
)

type testStruct struct {
	testName      string
	ruleName      string
	source        string
	expectedIsMet bool
}

var testDate = []testStruct{
	testStruct{
		testName: "TestRule1",
		ruleName: "rule1",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					public class SomePoorJavaClass {
						....
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule2NonArgsActivationFailurePath",
		ruleName: "rule2",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Activate
						public void activate() {
						}

						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule2NonArgsActivationHappyPath",
		ruleName: "rule2",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {

						@Activate
						@PostConstruct
						public void activate() {
						}

						....
					}
				`,
		expectedIsMet: true,
	},
	testStruct{
		testName: "TestRule2ArgsActivationFailurePath",
		ruleName: "rule2",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Activate
						public void activate(Map<String, Object> properties) {
						}

						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule2ArgsActivationHappyPath",
		ruleName: "rule2",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {

						public SomePoorJavaClass(){}

						@Autowired
						public SomePoorJavaClass(@Value("${xxxx}" String xxx)){
							activate(ImmutableMap.of(yyy, xxx));
						}
						
						@Activate
						public void activate(Map<String, Object> properties) {
						}

						...
					}
				`,
		expectedIsMet: true,
	},
	testStruct{
		testName: "TestRule3NonArgsDeactivationFailurePath",
		ruleName: "rule3",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Deactivate
						public void deactivate() {
						}

						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule3NonArgsDeactivationHappyPath",
		ruleName: "rule3",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Deactivate
						@PreDestroy
						public void deactivate() {
						}

						...
					}
				`,
		expectedIsMet: true,
	},
	testStruct{
		testName: "TestRule3ArgsDeactivationFailurePath",
		ruleName: "rule3",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Deactivate
						public void deactivate(Map<String, Object> properties) {
						}
						
						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule3ArgsDeactivationHappyPath",
		ruleName: "rule3",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {

						@PreDestroy
						public void tearDown(){
							deactivate(...)
						}
						
						@Deactivate
						public void deactivate(Map<String, Object> properties) {
						}
						
						...
					}
				`,
		expectedIsMet: true,
	},
	testStruct{
		testName: "TestRule4FailurePath",
		ruleName: "rule4",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Reference
						private AnotherPoorJavaClass random;

						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule4HappyPath",
		ruleName: "rule4",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Reference
						@Autowired
						private AnotherPoorJavaClass random;

						...
					}
				`,
		expectedIsMet: true,
	},
	testStruct{
		testName: "TestRule5FailurePath",
		ruleName: "rule5",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						@Autowired
						public SomePoorJavaClass(@Value("${xxxx}" String xxx)){
							activate(ImmutableMap.of(yyy, xxx));
						}
						
						@Activate
						public void activate(Map<String, Object> properties) {
						}

						...
					}
				`,
		expectedIsMet: false,
	},
	testStruct{
		testName: "TestRule5HappyPath",
		ruleName: "rule5",
		source: `
					...

					@Component(immediate=true, service = SomePoorJavaClass.class)
					@org.springframework.stereotype.Component
					public class SomePoorJavaClass {
						
						public SomePoorJavaClass(){}

						@Autowired
						public SomePoorJavaClass(@Value("${xxxx}" String xxx)){
							activate(ImmutableMap.of(yyy, xxx));
						}
						
						@Activate
						public void activate(Map<String, Object> properties) {
						}

						...
					}
				`,
		expectedIsMet: true,
	},
}

func TestRules(t *testing.T) {
	for _, data := range testDate {
		fmt.Printf("Running test case %v for rule: %v \n", data.testName, data.ruleName)
		rule := ruleRegistry[data.ruleName]
		isMet, err := rule.isMet(data.source)
		if err != nil {
			t.Errorf("Failed to run %v due to %v \n", data.testName, err)
			continue
		}

		if data.expectedIsMet != isMet {
			t.Errorf("Unexpected result for test case %v: %v \n", data.testName, isMet)
		}
	}
}

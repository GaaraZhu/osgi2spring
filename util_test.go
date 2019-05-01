package main

import (
	"fmt"
	"reflect"
	"testing"
)

type utilTestStruct struct {
	source             string
	expectedComponents []string
}

var utilTestData = []utilTestStruct{
	utilTestStruct{
		source: `
		....
		@Component
		public class Abc{}
				`,
		expectedComponents: []string{"Abc"},
	},
	utilTestStruct{
		source: `
		....
		@Component
		public abstract class Abc{}
				`,
		expectedComponents: []string{"Abc"},
	},
	utilTestStruct{
		source: `
		....
		@Component
		public class Abc<T extends D>{}
				`,
		expectedComponents: []string{"Abc"},
	},
	utilTestStruct{
		source: `
		....
		@Component
		public abstract class Abc <T extends D>{}
				`,
		expectedComponents: []string{"Abc"},
	},
	utilTestStruct{
		source: `
		....
		@Component
		public class Abc{
			@Component
			public static class AbcInnerA{}
			@Component
			public static class AbcInnerB{}
		}
				`,
		expectedComponents: []string{"Abc", "AbcInnerA", "AbcInnerB"},
	},
}

func TestComponentExtraction(t *testing.T) {
	fmt.Printf("Running test case to verify component extraction logic\n")
	for i, data := range utilTestData {
		actualComponents := extractOSGIComponents(data.source)
		if !reflect.DeepEqual(data.expectedComponents, actualComponents) {
			t.Errorf("Test case %v - Unexpected components %v expected: %v\n", i, actualComponents, data.expectedComponents)
		}
	}
}

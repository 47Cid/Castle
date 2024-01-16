package pod

import (
	"testing"
)

func TestInitFunc(t *testing.T) {
	t.Run("test case 1", func(t *testing.T) {
		// Call the function
		InitPods()

		// Check the results
		// For example, if initFunc() is supposed to set a global variable, you could check that:
		// if globalVar != expectedValue {
		// 	t.Errorf("got %v, want %v", globalVar, expectedValue)
		// }
	})
}

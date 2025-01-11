package assert

import (
	"strings"
	"testing"
)

// Notice how Equal() is a generic function? This means that we’ll be able to use it
// irrespective of what the type of the actual and expected values is. So long as both
// actual and expected have the same type and can be compared using the != operator
// (for example, they are both string values, or both int values).
func Equal[T comparable](t *testing.T, actual, expected T) {
	// Indicates to Go that this function is a test helper
	// Go test runner will report the filename and line number of the code which called the Equal() function
	// in the output.
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

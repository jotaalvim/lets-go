package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Helper marks the calling function as a test helper function
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want%v", actual, expected)
	}
}

func StringContains(t *testing.T, subString string, value string) {

	t.Helper()

	if !strings.Contains(value, subString) {
		t.Errorf("got %s; expected to contain%s", value, subString)
	}

}

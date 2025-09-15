package test

import (
	"testing"

	"github.com/k0kubun/pp"
)

func expect(t *testing.T, expected string, actual string) {
	if actual != expected {
		pp.Printf("Want: %v\nHave: %v\n", expected, actual)
		// t.Logf("Expected \"%s\", got \"%s\"", expected, actual)
		t.Fail()
	}
}

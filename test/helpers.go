package test

import (
	"testing"
)

func expect(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Logf("Expected \"%s\", got \"%s\"", expected, actual)
		t.Fail()
	}
}

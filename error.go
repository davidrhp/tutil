package tutil

import "testing"

func failOnError(t *testing.T, fn func() error) {
	if err := fn(); err != nil {
		t.Errorf("void function failed: %v", err)
	}
}

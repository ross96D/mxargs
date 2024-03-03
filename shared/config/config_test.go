package config

import (
	"testing"
)

func TestXxx(t *testing.T) {
	_, err := maxChars()
	if err != nil {
		t.Fatalf("err %s", err.Error())
	}
}

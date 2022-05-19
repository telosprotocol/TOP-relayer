package config

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	path := "./config.json"

	config, err := newConfig(path)
	if err != nil {
		t.Fatal("newConfig error:", err)
	}
	t.Log("config:", config)
	for _, chain := range config.Chains {
		t.Log("chain:", chain)
	}
}

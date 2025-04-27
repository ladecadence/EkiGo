package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	file := "../../testdata/testconfig.toml"
	conf, err := GetConfig(file)
	if err != nil {
		t.Errorf("Can't open config file: %v", err)
	}

	if conf.ID() != "MISSION" {
		t.Errorf("Expected MISSION, got %v", conf.ID())
	}

	if conf.LoraFreq() != 868500000 {
		t.Errorf("Expected 868500000, got %v", conf.LoraFreq())
	}
}

package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	file := "../../testdata/testconfig.toml"
	conf, err := GetConfig(file)
	if err != nil {
		t.Errorf("Can't read config file: %v", err)
	} else {
		if conf.ID() != "MISSION" {
			t.Errorf("Expected MISSION, got %v", conf.ID())
		}

		if conf.LoraFreq() != 868.5 {
			t.Errorf("Expected 868.5, got %v", conf.LoraFreq())
		}
	}
}

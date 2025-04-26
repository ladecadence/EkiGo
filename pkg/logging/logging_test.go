package logging

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	log, err := New("test")
	if err != nil {
		t.Errorf("Problem creating log file: %v", err)
	}
	err = log.Log(LogWarn, "Test message")
	if err != nil {
		t.Errorf("Problem adding data to log: %v", err)
	}

	data, err := os.ReadFile(log.Filename())
	if err != nil {
		panic(err)
	}

	fmt.Print(string(data))
	if !strings.Contains(string(data), "WARN") && !strings.Contains(string(data), "Test message") {
		t.Errorf("Can't find expected data in the log file")
	}

}

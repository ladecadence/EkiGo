package mcp3002

import (
	"fmt"
	"testing"
)

func TestADC(t *testing.T) {
	adc, err := New(1, 0)
	if err != nil {
		t.Errorf("Error starting ADC: %v", err)
	}

	value, err := adc.Read(0)
	if err != nil {
		t.Errorf("Error reading ADC: %v", err)
	}
	fmt.Printf("ADC value: %d\n", value)
}

package gps

import (
	"fmt"
	"testing"
)

func TestGPS(t *testing.T) {
	gps, err := New("/dev/serial0", 9600)
	if err != nil {
		t.Errorf("Error starting GPS: %v", err)
	}

	if gps != nil {
		defer gps.Close()
		fmt.Println("GPS Started")
		err = gps.Update()
		if err != nil {
			t.Errorf("Error updating GPS: %v", err)
		}
		h, m, s, err := gps.Time()
		if err != nil {
			t.Errorf("Error reading GPS time: %v", err)
		} else {
			fmt.Printf("GPS time: %v %v %v\n", h, m, s)
		}
	}
}

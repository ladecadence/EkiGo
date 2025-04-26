package led

import (
	"fmt"
	"testing"
	"time"
)

func TestLed(t *testing.T) {
	led, err := New(17)
	if err != nil {
		t.Errorf("Problem creating LED: %v", err)
	}

	// test LED
	fmt.Println("Normal blink")
	err = led.Blink()
	if err != nil {
		t.Errorf("Problem blinking LED: %v", err)
	}

	time.Sleep(time.Second * 2)
	// error
	fmt.Println("Error blink")
	err = led.BlinkError()
	if err != nil {
		t.Errorf("Problem error blinking LED: %v", err)
	}
}

package ms5607

import (
	"fmt"
	"testing"
)

func TestMS5607(t *testing.T) {
	ms, err := New(1, 0x77)
	if err != nil {
		t.Errorf("Problem starting MS5607: %v", err)
	}

	if ms != nil {
		err = ms.Update()
		if err != nil {
			t.Errorf("Problem updating data: %v", err)
		} else {
			fmt.Printf("MS5607 data: %fºC, %fmBar\n", ms.GetTemp(), ms.GetPres())
		}
	}
}

package rf95

import (
	"testing"
)

func TestRF95(t *testing.T) {
	rf, err := New(0, 0, 25, false)
	if err != nil {
		t.Errorf("Problem starting RF95: %v", err)
	}

	if rf != nil {
		err = rf.SetFrequency(868.5)
		if err != nil {
			t.Errorf("Problem setting frequency: %v", err)
		}
	}
}

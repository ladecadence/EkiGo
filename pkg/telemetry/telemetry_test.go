package telemetry

import (
	"fmt"
	"strings"
	"testing"
)

func TestTelemetry(t *testing.T) {
	telem := New("TEST", "Test telemetry message", "/")

	telem.Update(4332.944, "N", 539.783, "W", 0.0, 0.0, 0.0, 0, 0.0, 1019.5, 15.5, 5.4, 0)

	aprs := telem.AprsString()
	fmt.Println(aprs)

	if !strings.Contains(aprs, "P=1019.5") {
		t.Errorf("Problem with generated APRS string: %s", aprs)
	}

}

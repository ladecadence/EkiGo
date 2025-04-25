package picture

import (
	"testing"
)

func TestPicture(t *testing.T) {
	pic := New(0, "test", "./")

	err := pic.AddInfo("test.jpg", "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info tom image: %v", err)
	}
}

package picture

import (
	"fmt"
	"testing"
)

func TestPictureData(t *testing.T) {
	pic := New(0, "test", "./")

	err := pic.AddInfo("test.jpg", "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info tom image: %v", err)
	}

}

func TestPictureTake(t *testing.T) {
	pic := New(0, "test", "/home/pi/eki2/images")
	pic.UpdateName()
	pic.Capture()

	// now small picture
	pic2 := New(0, "small", "/home/pi/eki2/images")
	pic2.CaptureSmall("ssdv"+fmt.Sprintf("%d", pic2.Number), "640x480")

	err := pic.AddInfo(pic2.Filename, "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info tom image: %v", err)
	}

}

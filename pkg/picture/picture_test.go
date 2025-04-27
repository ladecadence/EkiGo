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
	pic := New(0, "test", "/home/pi/eki2/image/")
	pic.UpdateName()
	pic.Capture()
	fmt.Printf("Picture shot: %s\n", pic.Filename)

	// now small picture
	pic.CaptureSmall("ssdv", "640x480")

	err := pic.AddInfo(pic.Path+"ssdv", "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info tom image: %v", err)
	}

}

package picture

import (
	"fmt"
	"testing"
)

func TestPictureData(t *testing.T) {
	pic := New(0, "test", "./")

	err := pic.AddInfo("test.jpg", "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info to image: %v", err)
	}

}

func TestPictureTake(t *testing.T) {
	pic := New(0, "test", "/home/pi/eki2/image/")
	pic.UpdateName()
	err := pic.Capture(true)
	if err != nil {
		t.Errorf("Problem capturing image: %v", err)
	}
	fmt.Printf("Picture shot: %s\n", pic.Filename)

	// now small picture
	err = pic.CaptureSmall("ssdv.jpg", "640x480")
	if err != nil {
		t.Errorf("Problem capturing small image: %v", err)
	}
	fmt.Printf("Picture shot: %s\n", pic.Path+"ssdv.jpg")

	err = pic.AddInfo(pic.Path+"ssdv.jpg", "ID", "USBID", "Test image, message", "Data: 1234567890")
	if err != nil {
		t.Errorf("Problem adding info to image: %v", err)
	}

}

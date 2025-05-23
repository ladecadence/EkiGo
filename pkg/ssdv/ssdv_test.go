package ssdv

import (
	"fmt"
	"testing"
)

func TestSSDV(t *testing.T) {
	ss := New("/home/pi/eki2/image/ssdv.jpg", "/home/pi/eki2/image/ssdv/", "ssdv", "EKI2", 0)

	err := ss.Encode()
	if err != nil {
		t.Errorf("Error encoding SSDV image: %v", err)
	}

	packet, err := ss.GetPacket(100)
	if err != nil {
		t.Errorf("Error getting SSDV packet: %v", err)
	}
	if packet != nil {
		fmt.Printf("% x\n", packet)
	}
}

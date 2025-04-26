package ds18b20

import (
	"fmt"
	"testing"
)

func TestGetTemp(t *testing.T) {
	// do not init, use fake device file
	ds := DS18B20{
		Device: "../../testdata/ds18b20test.txt",
		Temp:   999.99,
	}

	temp, err := ds.Read()
	if err != nil {
		t.Errorf("Error reading temperature %v", err)
	}

	if temp != 19.937 {
		t.Errorf("Expecting 19.937, got %v", temp)
	}

	// real test
	ds.Init("28-031682a91bff")
	temp, err = ds.Read()
	if err != nil {
		t.Errorf("Error reading temperature %v", err)
	}
	fmt.Printf("Temperature: %v\n", temp)
}

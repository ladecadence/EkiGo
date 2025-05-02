package ssdv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	ssdvCommand = "ssdv"
)

type SSDV struct {
	imageFile  string
	id         string
	count      uint8
	fileName   string
	binaryName string
	Packets    uint64
}

func New(img string, path string, name string, id string, count uint8) SSDV {
	ss := SSDV{
		imageFile:  img,
		id:         id,
		count:      count,
		fileName:   img,
		binaryName: path + name + ".bin",
	}

	return ss
}

func (s *SSDV) Encode() error {
	cmd := exec.Command(ssdvCommand)
	cmd.Args = append(cmd.Args, "-e")
	cmd.Args = append(cmd.Args, "-c")
	cmd.Args = append(cmd.Args, s.id)
	cmd.Args = append(cmd.Args, "-i")
	cmd.Args = append(cmd.Args, fmt.Sprintf("%d", s.count))
	cmd.Args = append(cmd.Args, s.fileName)
	cmd.Args = append(cmd.Args, s.binaryName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("SSDV command: %v", cmd.Stderr)
		return err
	}
	// ssdv worked, get number of packets and return
	fi, err := os.Stat(s.binaryName)
	if err != nil {
		return err
	}
	// calculate number of packets of the file
	s.Packets = uint64(fi.Size()) / 256

	return nil
}

func (s *SSDV) GetPacket(packet uint64) ([]uint8, error) {
	if s.Packets == 0 {
		return nil, errors.New("No packets")
	}

	if packet > s.Packets {
		return nil, errors.New("Invalid packet")
	}

	// ok, try to get packet
	fi, err := os.Open(s.binaryName)
	if err != nil {
		return nil, errors.New("Can't open file")
	}
	defer fi.Close()

	// create a buffer to read
	buf := make([]byte, 255)
	_, err = fi.Seek((int64(packet)*256)+1, 0)
	b, err := fi.Read(buf)
	if err != nil || b < 255 {
		return nil, errors.New("Cant read data")
	}

	return buf, nil
}

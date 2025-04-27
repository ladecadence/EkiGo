package picture

import (
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

// Constants
const (
	raspistill         = "rpicam-still"
	TEXT_BIG   float32 = 20.0
	TEXT_SMALL float32 = 16.0
)

type Picture struct {
	Number   uint8
	Basename string
	Path     string
	Filename string
	captured bool
}

func New(num uint8, name string, path string) Picture {
	return Picture{
		Number:   num,
		Basename: name,
		Filename: path + name + string(num) + ".jpg",
		Path:     path,
		captured: false,
	}
}

func (p *Picture) UpdateName() {
	p.Filename = p.Path +
		p.Basename +
		"-" +
		time.Now().Format(time.RFC3339) +
		"-" +
		fmt.Sprintf("%d", p.Number) +
		".jpg"
}

func (p *Picture) Capture(rotate bool) error {
	p.UpdateName()
	var rot = ""
	if rotate == true {
		rot = "-r 180"
	}
	cmd := exec.Command(raspistill,
		"-t",
		"1000",
		rot,
		"-o",
		p.Filename,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	// if we manage to capture a picture,
	// increment filename number
	if p.Number == 255 {
		p.Number = 0
	} else {
		p.Number += 1
	}
	p.captured = true
	return nil
}

func (p *Picture) CaptureSmall(name string, res string) error {

	resolution := strings.Split(res, "x")

	cmd := exec.Command(raspistill,
		"-t",
		"1000",
		"--width",
		resolution[0],
		"--height",
		resolution[1],
		"-o",
		p.Path+name,
	)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *Picture) AddInfo(file string, id string, subid string, msg string, data string) error {
	datetime := time.Now().Format(time.RFC3339)

	// try to open image
	image, err := gg.LoadImage(file)
	if err != nil {
		return err
	}
	// context
	ct := gg.NewContextForImage(image)

	// font
	if err := ct.LoadFontFace("TerminusTTF-4.46.0.ttf", 20); err != nil {
		return err
	}
	// add texts
	ct.SetRGB(0, 0, 0)
	ct.DrawString(fmt.Sprintf("%s %s", id, subid), 10, 20)
	ct.SetRGB(1, 1, 1)
	ct.DrawString(fmt.Sprintf("%s %s", id, subid), 12, 22)

	// font
	if err := ct.LoadFontFace("TerminusTTF-4.46.0.ttf", 16); err != nil {
		return err
	}
	ct.SetRGB(0, 0, 0)
	ct.DrawString(msg, 10, 45)
	ct.SetRGB(1, 1, 1)
	ct.DrawString(msg, 11, 46)

	ct.SetRGB(0, 0, 0)
	ct.DrawString(datetime, 10, 65)
	ct.SetRGB(1, 1, 1)
	ct.DrawString(datetime, 11, 66)

	ct.SetRGB(0, 0, 0)
	ct.DrawString(data, 10, 80)
	ct.SetRGB(1, 1, 1)
	ct.DrawString(data, 11, 81)

	image = ct.Image()
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Specify the quality, between 0-100
	// Higher is better
	opt := jpeg.Options{
		Quality: 100,
	}
	err = jpeg.Encode(f, image, &opt)
	if err != nil {
		return err
	}

	return nil
}

package led

import (
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type LED interface {
	Blink() error
	BlinkError() error
}

type led struct {
	Pin   gpio.PinIO
	Ready bool
}

func New(pin uint8) (LED, error) {
	l := led{Ready: false}

	// get state TODO
	_, err := host.Init()
	if err != nil {
		return nil, err
	}
	// if err := rpio.Open(); err != nil {
	// 	return nil, err
	// }
	// Unmap gpio memory when done
	//defer rpio.Close()

	// configure
	l.Pin = gpioreg.ByName(fmt.Sprintf("%d", pin))
	l.Pin.Out(gpio.Low)
	l.Ready = true
	return &l, nil
}

func (l *led) Blink() error {
	if !l.Ready {
		return errors.New("LED GPIO Not initialized")
	}
	// if err := rpio.Open(); err != nil {
	// 	return err
	// }
	// Unmap gpio memory when done
	//defer rpio.Close()

	go func() {
		l.Pin.Out(gpio.High)
		time.Sleep(time.Millisecond)
		l.Pin.Out(gpio.Low)
	}()
	return nil
}

func (l *led) BlinkError() error {
	if !l.Ready {
		return errors.New("LED GPIO Not initialized")
	}
	// if err := rpio.Open(); err != nil {
	// 	return err
	// }
	// // Unmap gpio memory when done
	// defer rpio.Close()

	go func() {
		for range 5 {
			l.Pin.Out(gpio.High)
			time.Sleep(time.Millisecond * 5)
			l.Pin.Out(gpio.Low)
			time.Sleep(time.Millisecond * 5)
		}
		// keep on
		l.Pin.Out(gpio.High)
	}()
	return nil
}

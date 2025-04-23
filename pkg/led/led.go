package led

import (
	"errors"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type LED interface {
	Blink() error
	BlinkError() error
}

type led struct {
	Pin   rpio.Pin
	Ready bool
}

func New(pin int) (LED, error) {
	l := led{Ready: false}

	if err := rpio.Open(); err != nil {
		return nil, err
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	// configure
	l.Pin = rpio.Pin(pin)
	l.Pin.Output()
	l.Ready = true
	return &l, nil
}

func (l *led) Blink() error {
	if !l.Ready {
		return errors.New("LED GPIO Not initialized")
	}
	if err := rpio.Open(); err != nil {
		return err
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	go func() {
		l.Pin.High()
		time.Sleep(time.Millisecond)
		l.Pin.Low()
	}()
	return nil
}

func (l *led) BlinkError() error {
	if !l.Ready {
		return errors.New("LED GPIO Not initialized")
	}
	if err := rpio.Open(); err != nil {
		return err
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	go func() {
		for range 5 {
			l.Pin.High()
			time.Sleep(time.Millisecond)
			l.Pin.Low()
			time.Sleep(time.Millisecond)
		}
		// keep on
		l.Pin.High()
	}()
	return nil
}

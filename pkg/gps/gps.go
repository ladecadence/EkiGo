package gps

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	fieldTime  = 1
	fieldLat   = 2
	fieldNS    = 3
	fieldLon   = 4
	fieldEW    = 5
	fieldSats  = 7
	fieldAlt   = 9
	fieldSpeed = 7
	fieldHDG   = 8
	fieldDate  = 9

	minSats = 4
)

type GPS interface {
	Close() error
	Update() error
	Lat() float64
	NS() string
	Lon() float64
	EW() string
	Alt() float64
	Sats() int
	Hdg() float64
	Spd() float64
	Date() string
	Time() string
	Hms() (int, int, int, error)
	Dmy() (int, int, int, error)
}

type gps struct {
	time    string
	lat     float64
	ns      string
	lon     float64
	ew      string
	alt     float64
	sats    int
	hdg     float64
	spd     float64
	lineGGA string
	lineRMC string
	date    string
	port    serial.Port
}

func New(portFile string, speed int) (GPS, error) {
	// default values
	g := gps{
		time:    "",
		lat:     4332.944,
		ns:      "N",
		lon:     539.783,
		ew:      "W",
		alt:     0.0,
		hdg:     0.0,
		sats:    0,
		spd:     0.0,
		lineGGA: "",
		lineRMC: "",
		date:    "",
		//port:    port,
	}

	// prepare port
	mode := &serial.Mode{
		BaudRate: speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// options := serial.OpenOptions{
	// 	PortName:        port,
	// 	BaudRate:        speed,
	// 	DataBits:        8,
	// 	StopBits:        1,
	// 	MinimumReadSize: 1,
	// }

	// open port
	var err error
	g.port, err = serial.Open(portFile, mode)
	if err != nil {
		return nil, err
	}

	// var err error
	// g.port, err = serial.Open(options)
	// if err != nil {
	// 	return &g, err
	// }

	return &g, nil
}

func (g *gps) Close() error {
	err := g.port.Close()
	if err != nil {
		return err
	}
	// 	var err error
	// 	if g.port != nil {
	// 		err = g.port.Close()
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	return nil
}

func (g *gps) Update() error {
	// Serial buffer
	buf := make([]byte, 512)
	char := make([]byte, 1)

	// start parsing buffer
	g.lineGGA = ""
	for g.lineGGA == "" {
		n := 0
		tStart := time.Now()
		for range 512 {
			b, err := g.port.Read(char)
			if err != nil {
				return err
			}
			// if we read something
			if b > 0 {
				buf[n] = char[0]
				n += 1
			}
			// timeout
			if tNow := time.Now(); tNow.Sub(tStart) > (time.Second * 5) {
				return errors.New("Timeout reading serial port")
			}
		}

		// try to find GGA data
		data := string(buf[:])
		if strings.Contains(data, "$GNGGA") {
			fmt.Println("Found GGA")
			// cut to the start
			data = data[strings.Index(data, "$GNGGA"):]
			// cut to endline
			if strings.Contains(data, "\n") {
				data = data[:strings.Index(data, "\n")]
				g.lineGGA = data
			}
		}
	}

	// now RMC line
	g.lineRMC = ""
	for g.lineRMC == "" {
		n := 0
		tStart := time.Now()
		for range 512 {
			b, err := g.port.Read(char)
			if err != nil {
				return err
			}
			// if we read something
			if b > 0 {
				buf[n] = char[0]
				n += 1
			}
			// timeout
			if tNow := time.Now(); tNow.Sub(tStart) > (time.Second * 5) {
				return errors.New("Timeout reading serial port")
			}
		}
		// try to find RMC data
		data := string(buf[:])
		if strings.Contains(data, "$GNRMC") {
			fmt.Println("Found RMC")
			// cut to the start
			data = data[strings.Index(data, "$GNRMC"):]
			// cut to endline
			if strings.Contains(data, "\n") {
				data = data[:strings.Index(data, "\n")]
				g.lineRMC = data
			}
		}
	}

	// ok we have both lines, parse them
	ggaData := strings.Split(g.lineGGA, ",")
	rmcData := strings.Split(g.lineRMC, ",")

	// enough fields?
	if len(ggaData) >= 9 && len(rmcData) >= 8 {
		// good fix ?
		var err error
		g.sats, err = strconv.Atoi(ggaData[fieldSats])
		if err != nil {
			return err
		}
		if g.sats < minSats {
			// not enough sats, but perhaps we can parse time and date
			g.time = ggaData[fieldTime]
			g.date = rmcData[fieldDate]
			return errors.New("Not enough sats")
		}
		// ok parse elements if possible, if not provide default values
		g.time = ggaData[fieldTime]
		g.lat, err = strconv.ParseFloat(ggaData[fieldLat], 64)
		if err != nil {
			g.lat = 0.0
		}
		g.ns = ggaData[fieldNS]
		if g.ns == "" {
			g.ns = "N"
		}
		g.lon, err = strconv.ParseFloat(ggaData[fieldLon], 64)
		if err != nil {
			g.lon = 0.0
		}
		g.ew = ggaData[fieldEW]
		if g.ew == "" {
			g.ew = "W"
		}
		g.alt, err = strconv.ParseFloat(ggaData[fieldAlt], 64)
		if err != nil {
			g.alt = 0.0
		}
		g.spd, err = strconv.ParseFloat(rmcData[fieldSpeed], 64)
		if err != nil {
			g.spd = 0.0
		}
		g.hdg, err = strconv.ParseFloat(rmcData[fieldHDG], 64)
		if err != nil {
			g.hdg = 0.0
		}
		g.date = rmcData[fieldDate]

	} else {
		return errors.New("GPS parse error, not enough fields")
	}

	return nil
}

// getters
func (g *gps) Lat() float64 { return g.lat }
func (g *gps) NS() string   { return g.ns }
func (g *gps) Lon() float64 { return g.lon }
func (g *gps) EW() string   { return g.ew }
func (g *gps) Alt() float64 { return g.alt }
func (g *gps) Sats() int    { return g.sats }
func (g *gps) Hdg() float64 { return g.hdg }
func (g *gps) Spd() float64 { return g.spd }
func (g *gps) Date() string { return g.date }
func (g *gps) Time() string { return g.time }

func (g *gps) Hms() (int, int, int, error) {
	if len(g.time) >= 6 {
		hour, err := strconv.Atoi(g.time[0:2])
		minute, err := strconv.Atoi(g.time[2:4])
		second, err := strconv.ParseFloat(g.time[4:], 64)
		if err != nil {
			return 0, 0, 0, errors.New("GPS time parse fields error")
		}
		return hour, minute, int(second), nil
	} else {
		return 0, 0, 0, errors.New("GPS time parse error")
	}
}

func (g *gps) Dmy() (int, int, int, error) {
	if len(g.date) >= 6 {
		day, err := strconv.Atoi(g.date[0:2])
		month, err := strconv.Atoi(g.date[2:4])
		year, err := strconv.Atoi(g.date[4:])
		if err != nil {
			return 0, 0, 0, errors.New("GPS date parse fields error")
		}
		return day, month, year + 2000, nil
	} else {
		return 0, 0, 0, errors.New("GPS date parse error")
	}
}

func NmeaToDec(latlon float64) float64 {
	degrees := math.Trunc(latlon / 100.0)
	fraction := (latlon - (degrees * 100.0)) / 60.0

	return degrees + fraction
}

package telemetry

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Telemetry interface {
	Update(lat float64,
		ns string,
		lon float64,
		ew string,
		alt float64,
		hdg float64,
		spd float64,
		sats int,
		vbat float64,
		baro float64,
		tin float64,
		tout float64,
		hpwr bool)
	AprsString() string
	CsvString() string
}

type telemetry struct {
	id       string
	msg      string
	lat      float64
	ns       string
	lon      float64
	ew       string
	alt      float64
	hdg      float64
	spd      float64
	sats     int
	vbat     float64
	baro     float64
	tin      float64
	tout     float64
	arate    float64
	date     string
	time     string
	sep      string
	dateTime time.Time
	hpwr     bool
}

func New(i string, m string, s string) Telemetry {
	dt := time.Now().UTC()

	// default values
	return &telemetry{
		id:    i,
		msg:   m,
		sep:   s,
		lat:   0.0,
		ns:    "N",
		lon:   0.0,
		ew:    "W",
		alt:   0.0,
		hdg:   0.0,
		spd:   0.0,
		sats:  0,
		vbat:  0.0,
		baro:  0.0,
		tin:   0.0,
		tout:  0.0,
		arate: 0.0,
		date:  fmt.Sprintf("%02d-%02d-%d", dt.Day(), dt.Month(), dt.Year()),
		time:  fmt.Sprintf("%02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second()),
		hpwr:  false,
	}
}

func (t *telemetry) Update(
	lat float64,
	ns string,
	lon float64,
	ew string,
	alt float64,
	hdg float64,
	spd float64,
	sats int,
	vbat float64,
	baro float64,
	tin float64,
	tout float64,
	hpwr bool) {

	// save old altitude for ascension rate
	oldAlt := t.alt

	// update fields
	t.lat = lat
	t.ns = ns
	t.lon = lon
	t.ew = ew
	t.alt = alt
	t.hdg = hdg
	t.spd = spd
	t.sats = sats
	t.vbat = vbat
	t.baro = baro
	t.tin = tin
	t.tout = tout
	t.hpwr = hpwr

	// save old datetime
	oldDateTime := t.dateTime

	// update packet date
	t.dateTime = time.Now().UTC()
	t.date = fmt.Sprintf("%02d-%02d-%d", t.dateTime.Day(), t.dateTime.Month(), t.dateTime.Year())
	t.time = fmt.Sprintf("%02d:%02d:%02d", t.dateTime.Hour(), t.dateTime.Minute(), t.dateTime.Second())

	// calculate ascension rate
	deltaTime := t.dateTime.Sub(oldDateTime)
	if deltaTime.Milliseconds() != 0 {
		t.arate = (t.alt - oldAlt) /
			(float64(deltaTime.Milliseconds()) / 1000.0)
	} else {
		t.arate = 0.0
	}
}

func (t *telemetry) AprsString() string {
	// gen APRS coordinates
	coords := fmt.Sprintf(
		"%07.2f%s%s%08.2f%s",
		t.lat, t.ns, t.sep, t.lon, t.ew)

	// gen APRS string
	aprs := "$$"
	aprs += t.id
	aprs += "!"
	aprs += coords
	aprs += "O"
	aprs += fmt.Sprintf("%.1f", t.hdg)
	aprs += t.sep
	aprs += fmt.Sprintf("%.1f", t.spd)
	aprs += t.sep
	aprs += fmt.Sprintf("A=%.1f", t.alt)
	aprs += t.sep
	aprs += fmt.Sprintf("V=%.1f", t.vbat)
	aprs += t.sep
	aprs += fmt.Sprintf("P=%.1f", t.baro)
	aprs += t.sep
	aprs += fmt.Sprintf("TI=%.1f", t.tin)
	aprs += t.sep
	aprs += fmt.Sprintf("TO=%.1f", t.tout)
	aprs += t.sep
	aprs += t.date
	aprs += t.sep
	aprs += t.time
	aprs += t.sep
	aprs += fmt.Sprintf(
		"GPS=%09.6f%s,%010.6f%s",
		decLat(t.lat),
		t.ns,
		decLon(t.lon),
		t.ew)
	aprs += t.sep
	aprs += fmt.Sprintf("SATS=%d", t.sats)
	aprs += t.sep
	aprs += fmt.Sprintf("AR=%.1f", t.arate)
	aprs += t.sep
	aprs += strings.ReplaceAll(t.msg, "\n", " - ")
	aprs += func() string {
		if t.hpwr {
			return " - H"
		} else {
			return " - L"
		}
	}()
	aprs += "\n"

	return aprs
}

func (t *telemetry) CsvString() string {
	// gen CSV string
	csv := ""
	csv += t.date + ","
	csv += t.time + ","
	csv += fmt.Sprintf("%f", decLat(t.lat)) + ","
	csv += t.ns + ","
	csv += fmt.Sprintf("%f", decLon(t.lon)) + ","
	csv += t.ew + ","
	csv += fmt.Sprintf("%.1f", t.alt) + ","
	csv += fmt.Sprintf("%.2f", t.vbat) + ","
	csv += fmt.Sprintf("%.1f", t.tin) + ","
	csv += fmt.Sprintf("%.1f", t.tout) + ","
	csv += fmt.Sprintf("%.1f", t.baro) + ","
	csv += fmt.Sprintf("%.1f", t.hdg) + ","
	csv += fmt.Sprintf("%.1f", t.spd) + ","
	csv += fmt.Sprintf("%d", t.sats) + ","
	csv += fmt.Sprintf("%.1f", t.arate) + ","

	return csv
}

func decLat(lat float64) float64 {
	degrees := math.Trunc(lat / 100.0)
	fraction := (lat - (degrees * 100.0)) / 60.0
	return degrees + fraction
}

func decLon(lon float64) float64 {
	degrees := math.Trunc(lon / 100.0)
	fraction := (lon - (degrees * 100.0)) / 60.0
	return degrees + fraction
}

package mission

import (
	"fmt"
	"time"

	"github.com/ladecadence/EkiGo/pkg/config"
	"github.com/ladecadence/EkiGo/pkg/ds18b20"
	"github.com/ladecadence/EkiGo/pkg/gps"
	"github.com/ladecadence/EkiGo/pkg/led"
	"github.com/ladecadence/EkiGo/pkg/logging"
	"github.com/ladecadence/EkiGo/pkg/mcp3002"
	"github.com/ladecadence/EkiGo/pkg/ms5607"
	"github.com/ladecadence/EkiGo/pkg/picture"
	"github.com/ladecadence/EkiGo/pkg/rf95"
	"github.com/ladecadence/EkiGo/pkg/telemetry"
)

type Mission interface {
	Gps() gps.GPS
	Log() logging.Logging
	DataLog() logging.Logging
	UpdateTelemetry(config.Config) error
	SendTelemetry() error
	Telemetry() telemetry.Telemetry
	SetTimeGPS(int, int, int) error
}

type mission struct {
	log           logging.Logging
	dataLog       logging.Logging
	gps           gps.GPS
	led           led.LED
	adc           mcp3002.MCP3002
	baro          ms5607.MS5607
	temp_internal ds18b20.DS18B20
	temp_external ds18b20.DS18B20
	lora          rf95.RF95
	telem         telemetry.Telemetry
	pic           picture.Picture
	pwrSel        uint8
}

func New(conf config.Config) (Mission, error) {
	mission := mission{}

	// log
	var err error
	mission.log, err = logging.New(conf.PathMainDir() + conf.PathLogPrefix())
	if err != nil {
		return nil, err
	}
	mission.log.Log(logging.LogInfo, fmt.Sprintf("Mission %s starting...\n", conf.ID()))

	// datalog
	mission.dataLog, err = logging.New(conf.PathMainDir() + "datalog_")
	if err != nil {
		return nil, err
	}

	// gps
	mission.gps, err = gps.New(conf.GpsPort(), conf.GpsSpeed())
	if err != nil {
		return nil, err
	}

	// status led
	mission.led, err = led.New(conf.LedPin())
	if err != nil {
		return nil, err
	}

	// ADC and battery TODO: Battery enable pin
	mission.adc, err = mcp3002.New(conf.ADCCsPin(), conf.ADCChan())
	if err != nil {
		return nil, err
	}

	// barometer
	mission.baro, err = ms5607.New(conf.BaroI2CBus(), conf.BaroI2CAddr())
	if err != nil {
		return nil, err
	}

	// LoRa radio
	mission.lora, err = rf95.New(conf.LoraSPIChannel(), conf.LoraCSPin(), conf.LoraIntPin(), false)
	if err != nil {
		return nil, err
	}
	mission.lora.SetFrequency(conf.LoraFreq())

	// power selection
	// TODO read power selection pin
	mission.lora.SetTxPower(conf.LoraLowPwr())

	// telemetry
	mission.telem = telemetry.New(conf.ID(), conf.Msg(), conf.Separator())

	return &mission, nil
}

func (m *mission) Gps() gps.GPS {
	return m.gps
}

func (m *mission) Log() logging.Logging {
	return m.log
}

func (m *mission) DataLog() logging.Logging {
	return m.dataLog
}

func (m *mission) Telemetry() telemetry.Telemetry {
	return m.telem
}

func (m *mission) SetTimeGPS(hour, minute, second int) error {
	return nil
}

func (m *mission) UpdateTelemetry(conf config.Config) error {
	// Update sensor data
	// GPS
	err := m.gps.Update()
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData,
		fmt.Sprintf("%f%s, %f%s, Alt: %.1fm, Sats: %d, Date: %s, Time: %s",
			gps.NmeaToDec(m.gps.Lat()),
			m.gps.NS(),
			gps.NmeaToDec(m.gps.Lon()),
			m.gps.EW(),
			m.gps.Alt(),
			m.gps.Sats(),
			m.Gps().Date(),
			m.Gps().Time(),
		),
	)

	// baro
	err = m.baro.Update()
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData, fmt.Sprintf("BARO: %f", m.baro.GetPres()))

	// temperatures
	tin, err := m.temp_internal.Read()
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData, fmt.Sprintf("TIN: %.2f", tin))

	tout, err := m.temp_external.Read()
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData, fmt.Sprintf("TOUT: %.2f", tout))

	// Battery, enable reading, read ADC channel and make conversion
	// enable batt reading TODO
	// wait 1ms for current to stabilize
	time.Sleep(time.Millisecond * 1)
	adcBatt, err := m.adc.Read(conf.ADCVBatt())
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData, fmt.Sprintf("ADC0: %d", adcBatt))

	// disable batt reading TODO
	// convert to volts
	vBatt := conf.ADCVMult() * conf.ADCVDivider() * (float64(adcBatt) * 3.3 / 1023.0)
	m.log.Log(logging.LogData, fmt.Sprintf("VBATT: %.1f", vBatt))

	// Create telemetry packet
	m.Telemetry().Update(
		m.gps.Lat(),
		m.gps.NS(),
		m.gps.Lon(),
		m.gps.EW(),
		m.gps.Alt(),
		m.gps.Hdg(),
		m.gps.Spd(),
		m.gps.Sats(),
		vBatt,
		m.baro.GetPres(),
		tin,
		tout,
		m.pwrSel)

	fmt.Printf("APRS: %s\n", m.Telemetry().AprsString())

	return nil
}

func (m *mission) SendTelemetry() error {
	err := m.log.Log(logging.LogInfo, "Sending telemetry packet...")
	if err != nil {
		return err
	}
	err = m.lora.Send([]uint8(m.telem.AprsString()))
	if err != nil {
		return err
	}
	m.lora.WaitPacketSent()
	err = m.log.Log(logging.LogInfo, "Telemetry packet sent.")
	if err != nil {
		return err
	}
	err = m.led.Blink()
	if err != nil {
		return err
	}
	return nil
}

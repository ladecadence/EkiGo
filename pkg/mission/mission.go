package mission

import (
	"fmt"
	"time"

	"github.com/ladecadence/EkiGo/pkg/batt"
	"github.com/ladecadence/EkiGo/pkg/config"
	"github.com/ladecadence/EkiGo/pkg/ds18b20"
	"github.com/ladecadence/EkiGo/pkg/gps"
	"github.com/ladecadence/EkiGo/pkg/led"
	"github.com/ladecadence/EkiGo/pkg/logging"
	"github.com/ladecadence/EkiGo/pkg/mcp3002"
	"github.com/ladecadence/EkiGo/pkg/ms5607"
	"github.com/ladecadence/EkiGo/pkg/picture"
	"github.com/ladecadence/EkiGo/pkg/pwrsel"
	"github.com/ladecadence/EkiGo/pkg/rf95"
	"github.com/ladecadence/EkiGo/pkg/ssdv"
	"github.com/ladecadence/EkiGo/pkg/telemetry"
)

type Mission interface {
	Gps() gps.GPS
	Log() logging.Logging
	DataLog() logging.Logging
	UpdateTelemetry(config.Config) error
	SendTelemetry() error
	SendSSDV(config.Config) error
	Telemetry() telemetry.Telemetry
	SetTimeGPS(int, int, int) error
}

type mission struct {
	log           logging.Logging
	dataLog       logging.Logging
	gps           gps.GPS
	led           led.LED
	adc           mcp3002.MCP3002
	batt          batt.Batt
	baro          ms5607.MS5607
	temp_internal ds18b20.DS18B20
	temp_external ds18b20.DS18B20
	lora          rf95.RF95
	telem         telemetry.Telemetry
	pic           picture.Picture
	ssdv          ssdv.SSDV
	pwrSel        pwrsel.Pwrsel
}

func New(conf config.Config) (Mission, error) {
	mission := mission{}

	// log
	var err error
	mission.log, err = logging.New(conf.PathMainDir() + conf.PathLogPrefix())
	if err != nil {
		return nil, err
	}
	mission.log.Log(logging.LogInfo, fmt.Sprintf("Mission %s starting...", conf.ID()))

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

	// ADC and battery
	mission.adc, err = mcp3002.New(conf.ADCCsPin(), conf.ADCChan())
	if err != nil {
		return nil, err
	}

	mission.batt, err = batt.New(conf.BattEnablePin(), conf.ADCVMult(), conf.ADCVDivider())
	if err != nil {
		return nil, err
	}

	// barometer
	mission.baro, err = ms5607.New(conf.BaroI2CBus(), conf.BaroI2CAddr())
	if err != nil {
		return nil, err
	}

	// temperature sensors
	mission.temp_internal = ds18b20.DS18B20{}
	mission.temp_internal.Init(conf.TempInternalAddr())
	mission.temp_external = ds18b20.DS18B20{}
	mission.temp_external.Init(conf.TempInternalAddr())

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

	// picture
	mission.pic = picture.New(0, conf.ID(), conf.PathMainDir()+conf.PathImgDir())

	// ssdv
	mission.ssdv = ssdv.New(
		conf.PathMainDir()+conf.PathImgDir()+conf.SsdvName(),
		conf.PathMainDir()+conf.PathImgDir(),
		conf.SsdvName(),
		conf.ID(),
		mission.pic.Number,
	)

	// pwr selection pin
	mission.pwrSel, err = pwrsel.New(conf.PwrPin())
	if err != nil {
		return nil, err
	}

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

	// Battery
	vBatt, err := m.batt.Read(m.adc, uint8(conf.ADCChan()))
	if err != nil {
		return err
	}
	m.log.Log(logging.LogData, fmt.Sprintf("VBATT: %.1f", vBatt))

	pwrSel := m.pwrSel.Read()

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
		pwrSel)

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

func (m *mission) SendSSDV(conf config.Config) error {
	err := m.pic.Capture(true)
	if err != nil {
		m.log.Log(logging.LogError, fmt.Sprintf("Error taking picture: %v", err))
		return err
	}
	m.log.Log(logging.LogInfo, fmt.Sprintf("Picture shot: %s", m.pic.Filename))
	// Take SSDV picture
	err = m.pic.CaptureSmall(conf.SsdvName(), conf.SsdvSize(), true)
	if err != nil {
		m.log.Log(logging.LogError, fmt.Sprintf("Error taking SSDV picture: %v", err))
		return err
	}
	m.log.Log(logging.LogInfo, fmt.Sprintf("SSDV picture shot: %s", conf.SsdvName()))

	// Encode SSDV picture
	// Add info to image
	err = m.pic.AddInfo(conf.PathMainDir()+conf.PathImgDir()+conf.SsdvName(),
		conf.ID(),
		conf.SubID(),
		conf.Msg(),
		fmt.Sprintf("%f%s, %f%s, %.1fm",
			gps.NmeaToDec(m.gps.Lat()),
			m.gps.NS(),
			gps.NmeaToDec(m.gps.Lon()),
			m.gps.EW(),
			m.gps.Alt(),
		),
	)
	if err != nil {
		m.log.Log(logging.LogError, fmt.Sprintf("Error adding info to SSDV picture: %v", err))
		return err
	}
	m.log.Log(logging.LogInfo, "SSDV info added")

	m.ssdv = ssdv.New(
		conf.PathMainDir()+conf.PathImgDir()+conf.SsdvName(),
		conf.PathMainDir()+conf.PathImgDir(),
		conf.SsdvName(),
		conf.ID(),
		m.pic.Number,
	)

	// launch SSDV to create bin SSDV img
	err = m.ssdv.Encode()
	if err != nil {
		m.log.Log(logging.LogError, fmt.Sprintf("Error encoding SSDV binary file: %v", err))
		return err
	}

	// send it
	err = m.log.Log(logging.LogInfo, "Sending SSDV picture...")
	if err != nil {
		return err
	}
	lastTime := time.Now()
	for i := range m.ssdv.Packets {
		packet, err := m.ssdv.GetPacket(i)
		if err != nil {
			return err
		}
		err = m.lora.Send(packet)
		if err != nil {
			return err
		}
		m.lora.WaitPacketSent()
		err = m.led.Blink()
		if err != nil {
			return err
		}
		err = m.log.Log(logging.LogInfo, fmt.Sprintf("SSDV sent packet %d.", i))
		if err != nil {
			return err
		}

		// check if we need to send telemetry between image packets
		if timeDiff := time.Now().Sub(lastTime); timeDiff > time.Second*time.Duration(conf.PacketDelay()) {
			err := m.UpdateTelemetry(conf)
			if err != nil {
				return err
			}
			err = m.SendTelemetry()
			if err != nil {
				return err
			}
			lastTime = time.Now()
		}

		// wait a bit between packets for decoding on the client
		time.Sleep(time.Millisecond * 100)
	}

	err = m.log.Log(logging.LogInfo, fmt.Sprintf("SSDV image, %d packets sent.", m.ssdv.Packets))
	if err != nil {
		return err
	}

	return nil
}

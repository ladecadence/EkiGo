package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ladecadence/EkiGo/pkg/config"
	"github.com/ladecadence/EkiGo/pkg/logging"
	"github.com/ladecadence/EkiGo/pkg/mission"
)

func main() {
	// command line flags
	configFile := flag.String("c", "config.toml", "Config file")
	flag.Parse()

	// try to load config file, if not, create a defualt configuration file
	// at standard config directory
	conf, err := config.GetConfig(*configFile)
	if err != nil {
		file, err := config.CreateDefaultConfig()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Can't open configuration, default configuration file created at %s\n", file)
		os.Exit(1)
	}

	// now test that configuration is not the default one
	if conf.ID() == "" || conf.SubID() == "" || conf.Msg() == "" ||
		conf.Separator() == "" || conf.PathMainDir() == "" ||
		conf.GpsPort() == "" {
		fmt.Println("Please edit the configuration file.")
		os.Exit(1)
	}

	// create mission and configure it
	mission, err := mission.New(conf)
	if err != nil {
		panic(err)
	}

	// Ok, now get time from GPS and update system time
	err = mission.Gps().Update()
	if err != nil {
		mission.Log().Log(logging.LogError, fmt.Sprintf("Error updating GPS: %v", err))
	}
	h, m, s, err := mission.Gps().Hms()
	if err != nil {
		mission.Log().Log(logging.LogError, fmt.Sprintf("Error getting GPS time:: %v", err))
	} else {
		err = mission.SetTimeGPS(h, m, s)
		if err != nil {
			mission.Log().Log(logging.LogError, fmt.Sprintf("Error setting system time: %v", err))
		}
	}

	///////// MAIN LOOP /////////
	for {
		// send Telemetry
		for range conf.PacketRepeat() {
			// check for commands TODO

			// send telemetry
			err := mission.UpdateTelemetry(conf)
			if err != nil {
				mission.Log().Log(logging.LogError, fmt.Sprintf("Error updating telemetry: %v", err))
			}
			err = mission.SendTelemetry()
			if err != nil {
				mission.Log().Log(logging.LogError, fmt.Sprintf("Problem sending telemetry: %v", err))
			}

			// write datalog
			err = mission.DataLog().Log(logging.LogClean, mission.Telemetry().CsvString())
			if err != nil {
				mission.Log().Log(logging.LogError, fmt.Sprintf("Problem writing datalog: %v", err))
			}

			time.Sleep(time.Duration(conf.PacketDelay()) * time.Second)
		}

		// send SSDV
		err = mission.SendSSDV(conf)
		if err != nil {
			mission.Log().Log(logging.LogError, fmt.Sprintf("Problem sending SSDV: %v", err))
		}
		time.Sleep(time.Duration(conf.PacketDelay()) * time.Second)
	}

}

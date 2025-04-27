package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ladecadence/EkiGo/pkg/config"
	"github.com/ladecadence/EkiGo/pkg/logging"
	"github.com/ladecadence/EkiGo/pkg/mission"
)

func main() {
	// command line flags
	configFile := flag.String("c", "config.toml", "Config file")
	flag.Parse()

	conf, err := config.GetConfig(*configFile)
	if err != nil {
		panic(err)
	}

	// now test that configuration is not the default one
	if conf.ID() == "" || conf.SubID() == "" || conf.Msg() == "" ||
		conf.Separator() == "" || conf.PathMainDir() == "" ||
		conf.GpsPort() == "" {
		fmt.Println("Please edit the configuration file.")
		panic("")
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
				fmt.Printf("Error updating telemetry: %v\n", err)
			}
			fmt.Println("Updated telemetry")
			mission.SendTelemetry()

			// write datalog
			mission.DataLog().Log(logging.LogClean, mission.Telemetry().CsvString())

			time.Sleep(time.Duration(conf.PacketDelay()) * time.Second)
		}

		// send SSDV
		// TODO
	}

}

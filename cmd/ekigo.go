package main

import (
	"flag"
	"fmt"

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

	fmt.Printf("%v\n", conf.ID())

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

}

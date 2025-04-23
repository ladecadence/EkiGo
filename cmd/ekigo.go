package main

import (
	"flag"
	"log"

	"github.com/ladecadence/EkiGo/pkg/config"
)

func main() {
	// command line flags
	configFile := flag.String("c", "config.toml", "Config file")
	flag.Parse()

	conf := config.GetConfig(*configFile)
	log.Printf("Mission %s starting...\n", conf.ID)
}

package main

import (
	"flag"
	"log"

	"cargoship/cmd/configurations"
)

func main() {
	// command line flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "extractor_config.yaml", "Path to configuration yaml")
	flag.Parse()

	// read script configuration
	configs, err := configurations.ReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}
	// read ftp times state
	times, err := configurations.ReadTimes(configs.TimesPath)
	if err != nil {
		log.Panic(err)
	}
	// defer update/write ftp time state file
	defer configurations.WriteTimes(&times, configs.TimesPath)

	extractFiles(*configs, &times)
}

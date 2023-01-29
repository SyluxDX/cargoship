package main

import (
	"flag"
	"fmt"
	"log"

	"cargoship/loader/cmd/configurations"
)

func main() {
	fmt.Println("Hello World")

	configurations.ReadConfig("loader_config.yaml")
	// command line flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "loader_config.yaml", "Path to configuration yaml")
	flag.Parse()

	// read script configuration
	configs, err := configurations.ReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configs)

	// // read ftp times state
	// times, err := configurations.ReadTimes(configs.TimesPath)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// // defer update/write ftp time state file
	// defer configurations.WriteTimes(&times, configs.TimesPath)
}

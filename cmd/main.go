package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"cargoship/cmd/configurations"

	"github.com/jlaffaye/ftp"
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

	for _, server := range configs.Ftps {
		log.Printf("Connect to server: %s\n", server.Name)
		// create connection to ftp
		ftpUrl := fmt.Sprintf("%s:%d", server.Hostname, server.Port)
		conn, err := ftp.Dial(ftpUrl, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Fatal(err)
		}
		// login
		err = conn.Login(server.User, server.Pass)
		if err != nil {
			log.Fatal(err)
		}
		// service loop
		for _, service := range server.Services {
			if service.Mode == "import" {
				extractFiles(server.Name, conn, service, &times)
			}
		}
	}
}

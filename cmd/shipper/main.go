package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"cargoship/cmd/shipper/transport"
	"cargoship/internal/configurations"
	"cargoship/internal/logging"

	"github.com/jlaffaye/ftp"
)

var (
	scriptLogger logging.Logger
	filesLogger  logging.Logger
)

func main() {
	// command line flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "shipper_config.yaml", "Path to configuration yaml")
	flag.Parse()

	// read script configuration
	configs, err := configurations.ShipperReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}
	// start loggers
	scriptLogger.Init(configs.Log.Script, configs.Log2Console)
	filesLogger.Init(configs.Log.Files, configs.Log2Console)
	defer scriptLogger.Close()
	defer filesLogger.Close()

	// read ftp times state
	times, err := configurations.ShipperReadTimes(configs.TimesPath)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}

	// defer update/write ftp time state file
	defer configurations.ShipperWriteTimes(&times, configs.TimesPath)

	for _, server := range configs.Ftps {
		scriptLogger.LogInfo(fmt.Sprintf("Connect to server: %s\n", server.Name))
		// create connection to ftp
		ftpUrl := fmt.Sprintf("%s:%d", server.Hostname, server.Port)
		conn, err := ftp.Dial(ftpUrl, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			scriptLogger.LogError(err.Error())
			panic(err)
		}
		// login
		err = conn.Login(server.User, server.Pass)
		if err != nil {
			scriptLogger.LogError(err.Error())
			panic(err)
		}
		// service loop
		for _, service := range server.Services {
			if service.Mode == "import" {
				transport.DownloadFiles(server.Name, conn, service, &times, scriptLogger, filesLogger)
			} else if service.Mode == "export" {
				transport.UploadFiles(server.Name, conn, service, &times, scriptLogger, filesLogger)
			} else {
				scriptLogger.LogWarn(
					fmt.Sprintf("ERROR Unknown mode, %s, on service %s.\n", service.Mode, service.Name),
				)
			}
		}
	}
}

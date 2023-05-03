package main

import (
	"flag"
	"fmt"
	"log"

	"cargoship/loader/cmd/configurations"
	"cargoship/loader/cmd/files"
	"cargoship/loader/cmd/logging"
)

// func debug_yaml(data interface{}) {
// 	out, _ := yaml.Marshal(data)
// 	fmt.Println(string(out))
// }

var (
	scriptLogger logging.Logger
	filesLogger  logging.Logger
)

func main() {
	// command line flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "loader_config.yaml", "Path to configuration yaml")
	flag.Parse()

	// read script configuration
	configs, err := configurations.ReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}

	// start loggers
	scriptLogger.Init(configs.Log.Script, configs.Log2Console)
	filesLogger.Init(configs.Log.Files, configs.Log2Console)
	defer scriptLogger.Close()
	defer filesLogger.Close()

	// read ftp times state
	times, err := configurations.ReadTimes(configs.TimesPath)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	// defer update/write ftp time state file
	defer times.WriteTimes(configs.TimesPath)

	// main loop
	for _, service := range configs.Services {
		if service.Mode == "compress" {
			fmt.Println("Compress mode", service.Name)
		} else if service.Mode == "cleaner" {
			files.CleanFiles(service, times, scriptLogger, filesLogger)
		} else {
			scriptLogger.LogWarn(
				fmt.Sprintf("ERROR Unknown mode, %s, on service %s.\n", service.Mode, service.Name),
			)
		}
		// scriptLogger.LogInfo(service.Name)

		// debug_yaml(service)
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cargoship/cmd/configurations"

	"github.com/jlaffaye/ftp"
	"gopkg.in/yaml.v3"
)

func checkFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(folderPath, 0755)
	}
}

func extractFiles(config configurations.ExtractorConfig, times configurations.FileTimes) {
	for _, server := range config.Ftps {
		fmt.Println(server.Name)
		for _, service := range server.Services {
			// check folder
			checkFolder(service.Dst)

			fmt.Println("", service.Name)
		}
	}
}

func test_ftp(config configurations.ExtractorConfig) {
	ftpServer := config.Ftps[0]

	ftpUrl := fmt.Sprintf("%s:%d", ftpServer.Hostname, ftpServer.Port)
	conn, err := ftp.Dial(ftpUrl, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Login(ftpServer.User, ftpServer.Pass)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := conn.List(ftpServer.Services[0].Src)
	if err != nil {
		log.Panicln(err)
	}
	for _, file := range entries {
		jdata, _ := json.MarshalIndent(file, "", "")
		fmt.Println(string(jdata))
		fmt.Printf("%T\n", file.Time)
	}

	if err := conn.Quit(); err != nil {
		log.Fatal(err)
	}
}

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

	// configurations.UpsertTimes(&times, "newftp", "bananas", time.Now().UTC().Format("2006-01-02T15:04:05"))
	aux, _ := yaml.Marshal(times)
	log.Println(string(aux))

	configurations.UpsertTimes(&times, "azure", "bannas", time.Now().UTC())

	aux, _ = yaml.Marshal(times)
	log.Println(string(aux))
	configurations.UpsertTimes(&times, "ftp", "hourlyprocess", time.Now().UTC())

	aux, _ = yaml.Marshal(times)
	log.Println(string(aux))

	// test_ftp(*configs)
	// extractFiles(*configs)
}

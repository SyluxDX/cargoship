package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

func listDirectory(conn *ftp.ServerConn, source string, prefix string, extension string) ([]*ftp.Entry, error) {
	var outputList []*ftp.Entry

	entries, err := conn.List(source)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name, prefix) && strings.HasSuffix(entry.Name, extension) {
			// entry.Name = fmt.Sprintf("%s/%s", source, entry.Name)
			outputList = append(outputList, entry)
		}
	}
	return outputList, nil
}

func dateFilterDirectory(entries []*ftp.Entry, lastTime time.Time) []*ftp.Entry {
	var outputList []*ftp.Entry

	for _, entry := range entries {
		if entry.Time.After(lastTime) {
			outputList = append(outputList, entry)
		}
	}
	return outputList
}

func fileDownload(conn *ftp.ServerConn, destination string, entry *ftp.Entry) error {
	remoteReader, err := conn.Retr(entry.Name)
	if err != nil {
		return err
	}
	defer remoteReader.Close()

	// create local writer
	localWriter, err := os.OpenFile(
		destination,
		// fmt.Sprintf("%s/%s", destination, entry.Name),
		// os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer localWriter.Close()

	sizeWritten, err := io.Copy(localWriter, remoteReader)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Donwloaded file %s (size %d), written %d\n", entry.Name, entry.Size, sizeWritten)
	return nil
}

func downloadFiles(conn *ftp.ServerConn, source string, destination string, entries []*ftp.Entry) (time.Time, error) {
	// move to source folder
	var lastFileTime time.Time

	err := conn.ChangeDir(source)
	if err != nil {
		return lastFileTime, err
	}
	// download files
	for _, entry := range entries {
		err := fileDownload(conn, fmt.Sprintf("%s/%s", destination, entry.Name), entry)
		if err != nil {
			return lastFileTime, err
		}
		// update
		lastFileTime = entry.Time
	}
	return lastFileTime, nil
}

func extractFiles(config configurations.ExtractorConfig, times []configurations.FileTimes) {
	for _, server := range config.Ftps {
		fmt.Println(server.Name)
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
			fmt.Println("", service.Name)
			// check folder
			checkFolder(service.Dst)
			// get last file time
			fileTime := configurations.GetTimes(times, server.Name, service.Name)
			fmt.Println(fileTime)
			// list files in directory
			entries, err := listDirectory(conn, service.Src, service.Prefix, service.Extension)
			if err != nil {
				log.Panic(err)
			}
			entries = dateFilterDirectory(entries, fileTime)
			for _, entry := range entries {
				fmt.Println(entry.Name)
			}
			// check if there are any files to download
			if len(entries) == 0 {
				return
			}

			dir, _ := conn.CurrentDir()
			fmt.Println(dir)
			// donwload files
			lastTime, err := downloadFiles(conn, service.Src, service.Dst, entries)
			if err != nil {
				log.Panic(err)
			}
			// update last downloaded time
			configurations.UpsertTimes(&times, server.Name, service.Name, lastTime)
			// fileTime := configurations.GetTimes(times, server.Name, service.Name)
			// fmt.Println(fileTime)
		}
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

	// configurations.UpsertTimes(&times, "azure", "bannas", time.Now().UTC())

	// aux, _ = yaml.Marshal(times)
	// log.Println(string(aux))
	// configurations.UpsertTimes(&times, "ftp", "hourlyprocess", time.Now().UTC())

	// aux, _ = yaml.Marshal(times)
	// log.Println(string(aux))

	// test_ftp(*configs)
	extractFiles(*configs, times)
	aux, _ = yaml.Marshal(times)
	log.Println(string(aux))
}

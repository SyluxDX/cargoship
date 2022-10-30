package configurations

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type configTimes struct {
	Ftp     string `yaml:"ftp"`
	Service string `yaml:"service"`
	Time    string `yaml:"time"`
}

type FileTimes struct {
	Ftp     string
	Service string
	Time    time.Time
}

func ReadTimes(filepath string) ([]FileTimes, error) {
	// create empty readConfig
	//var readConfig FileTimes
	var readConfig []configTimes

	// check if file exits
	if _, err := os.Stat(filepath); err == nil {
		//read file
		fdata, err := os.ReadFile((filepath))
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(fdata, &readConfig)
		if err != nil {
			return nil, err
		}

	}

	// convert string into time object
	var config []FileTimes
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config = append(config, FileTimes{
			Ftp:     serviceService.Ftp,
			Service: serviceService.Service,
			Time:    timestamp,
		})
	}

	return config, nil
}

// func GetTimes(ftpTimes []FileTimes, server string, service string) {
// 	for _, time := range ftpTimes {
// 		if time.Ftp == server && time.Service {

// 		}
// 	}
// }

func UpsertTimes(config *[]FileTimes, ftp string, service string, times time.Time) {
	for idx, filetimes := range *config {
		if filetimes.Ftp == ftp && filetimes.Service == service {
			(*config)[idx].Time = times
			return
		}
	}
	// new ftp times
	*config = append(*config, FileTimes{Ftp: ftp, Service: service, Time: times})
}

func WriteTimes(ftpTimes *[]FileTimes, filepath string) error {
	fileData := make([]configTimes, len(*ftpTimes))
	for idx, stamp := range *ftpTimes {
		fileData[idx] = configTimes{
			Ftp:     stamp.Ftp,
			Service: stamp.Service,
			Time:    stamp.Time.Format(time.RFC3339),
		}
	}

	// marshall times -> struct to string
	data, err := yaml.Marshal(fileData)
	if err != nil {
		return err
	}
	// write to file
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

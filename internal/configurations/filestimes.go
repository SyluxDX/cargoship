package configurations

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Loader structures
type configTimes struct {
	Service string `yaml:"service"`
	Mode    string `yaml:"mode"`
	Time    string `yaml:"time"`
}

// custom type used for receiver
type filetime struct {
	Service string
	Mode    string
	Time    time.Time
}
type FileTimes struct {
	Times []filetime
}

// Shipper structures
type shipperConfigTimes struct {
	Ftp     string `yaml:"ftp"`
	Mode    string `yaml:"mode"`
	Service string `yaml:"service"`
	Time    string `yaml:"time"`
}

type ShipperFileTimes struct {
	Ftp     string
	Mode    string
	Service string
	Time    time.Time
}

func ShipperReadTimes(filepath string) ([]ShipperFileTimes, error) {
	// create empty readConfig
	var readConfig []shipperConfigTimes

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
	var config []ShipperFileTimes
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config = append(config, ShipperFileTimes{
			Ftp:     serviceService.Ftp,
			Mode:    serviceService.Mode,
			Service: serviceService.Service,
			Time:    timestamp,
		})
	}

	return config, nil
}

//// change these functions to use receiver

func ShipperGetTimes(ftpTimes []ShipperFileTimes, server string, mode string, service string) time.Time {
	for _, elm := range ftpTimes {
		if elm.Ftp == server && elm.Service == service && elm.Mode == mode {
			return elm.Time
		}
	}
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

func ShipperUpsertTimes(config *[]ShipperFileTimes, ftp string, mode string, service string, times time.Time) {
	for idx, filetimes := range *config {
		if filetimes.Ftp == ftp && filetimes.Mode == mode && filetimes.Service == service {
			(*config)[idx].Time = times
			return
		}
	}
	// new ftp times
	*config = append(*config, ShipperFileTimes{Ftp: ftp, Mode: mode, Service: service, Time: times})
}

func ShipperWriteTimes(ftpTimes *[]ShipperFileTimes, filepath string) error {
	fileData := make([]shipperConfigTimes, len(*ftpTimes))
	for idx, stamp := range *ftpTimes {
		fileData[idx] = shipperConfigTimes{
			Ftp:     stamp.Ftp,
			Mode:    stamp.Mode,
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

//// end changes

func LoaderReadTimes(filepath string) (*FileTimes, error) {
	// create empty readConfig
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
	var config FileTimes
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config.Times = append(config.Times, filetime{
			Service: serviceService.Service,
			Mode:    serviceService.Mode,
			Time:    timestamp,
		})
	}

	return &config, nil
}

func (config *FileTimes) GetTimes(service, mode string) time.Time {
	for _, elm := range config.Times {
		if elm.Service == service && elm.Mode == mode {
			return elm.Time
		}
	}
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

func (config *FileTimes) UpsertTimes(service string, mode string, times time.Time) {
	for idx, filetimes := range config.Times {
		if filetimes.Mode == mode && filetimes.Service == service {
			config.Times[idx].Time = times
			return
		}
	}
	// new files times
	config.Times = append(config.Times, filetime{Service: service, Mode: mode, Time: times})
}

func (config *FileTimes) WriteTimes(filepath string) error {
	fileData := make([]configTimes, len(config.Times))
	for idx, stamp := range config.Times {
		fileData[idx] = configTimes{
			Service: stamp.Service,
			Mode:    stamp.Mode,
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

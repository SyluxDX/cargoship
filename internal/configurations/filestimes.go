package configurations

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

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

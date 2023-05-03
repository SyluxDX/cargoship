package configurations

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Shipper structures
type shipperConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"timesFilePath"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Ftps []struct {
		Name     string `yaml:"name"`
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		User     string `yaml:"username"`
		Pass     string `yaml:"password"`
		Protocol string `yaml:"protocol"`
	} `yaml:"ftps"`
	Services []struct {
		Name      string   `yaml:"name"`
		Ftp       []string `yaml:"ftpConfig"`
		Mode      string   `yaml:"mode"`
		Src       string   `yaml:"sourceFolder"`
		Dst       string   `yaml:"destinationFolder"`
		Prefix    string   `yaml:"filePrefix"`
		Extension string   `yaml:"fileExtension"`
		History   string   `yaml:"historyFolder"`
		MaxTime   int      `yaml:"maxTime"`
		Window    int      `yaml:"windowLimit"`
	} `yaml:"services"`
}
type ShipperService struct {
	Name      string
	Mode      string
	Src       string
	Dst       string
	Prefix    string
	Extension string
	History   string
	MaxTime   int
	Window    int
}
type FtpConfig struct {
	Name     string
	Hostname string
	Port     int
	User     string
	Pass     string
	Protocol string
	Services []ShipperService
}
type ShipperConfig struct {
	Log2Console bool
	TimesPath   string
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Ftps []FtpConfig
}

// Loader structures
type LoaderService struct {
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	Src       string `yaml:"sourceFolder"`
	Dst       string `yaml:"destinationFolder"`
	Prefix    string `yaml:"filePrefix"`
	Extension string `yaml:"fileExtension"`
	MaxTime   int    `yaml:"maxTime"`
	Window    int    `yaml:"windowLimit"`
}

type LoaderConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"timesFilePath"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Services []LoaderService `yaml:"services"`
}

func replaceDatePlaceholder(filename string) string {
	start := strings.Index(filename, "{")
	end := strings.Index(filename, "}")
	if start != -1 && end != -1 {
		dateFormat := filename[start+1 : end]
		return fmt.Sprintf(
			"%s%s%s",
			filename[:start],
			time.Now().UTC().Format(dateFormat),
			filename[end+1:],
		)
	}
	return filename
}

func LoaderReadConfig(filepath string) (*LoaderConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var loader LoaderConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &loader)
	if err != nil {
		return nil, err
	}

	loader.Log.Script = replaceDatePlaceholder(loader.Log.Script)
	loader.Log.Files = replaceDatePlaceholder(loader.Log.Files)

	return &loader, nil
}

func ShipperReadConfig(filepath string) (*ShipperConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config shipperConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	// process read config to match/duplicate service and server
	var newConfig = &ShipperConfig{
		Log2Console: config.Log2Console,
		TimesPath:   config.TimesPath,
		Log: struct {
			Script string `yaml:"script"`
			Files  string `yaml:"files"`
		}{
			Script: replaceDatePlaceholder(config.Log.Script),
			Files:  replaceDatePlaceholder(config.Log.Files),
		},
	}

	// create mapping service index
	ftpIndex := make(map[string]int, len(config.Ftps))
	for idx, ftp := range config.Ftps {
		ftpIndex[ftp.Name] = idx
		newConfig.Ftps = append(newConfig.Ftps, FtpConfig{
			Name:     ftp.Name,
			Hostname: ftp.Hostname,
			Port:     ftp.Port,
			User:     ftp.User,
			Pass:     ftp.Pass,
			Protocol: ftp.Protocol,
		})
	}

	// process serving matching to servers
	for _, service := range config.Services {
		match := ShipperService{
			Name:      service.Name,
			Mode:      service.Mode,
			Src:       service.Src,
			Dst:       service.Dst,
			Prefix:    service.Prefix,
			Extension: service.Extension,
			History:   service.History,
			MaxTime:   service.MaxTime,
			Window:    service.Window,
		}
		for _, ftpName := range service.Ftp {
			if idx, ok := ftpIndex[ftpName]; ok {
				newConfig.Ftps[idx].Services = append(newConfig.Ftps[idx].Services, match)
			}
		}
	}
	return newConfig, nil
}

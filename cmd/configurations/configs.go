package configurations

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type fileConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"timesFilePath"`
	Log         struct {
		Folder   string `yaml:"folder"`
		Filename string `yaml:"filename"`
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

type LogConfig struct {
	Folder   string
	Filename string
}

type ServiceConfig struct {
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
	Services []ServiceConfig
}
type ExtractorConfig struct {
	Log2Console bool
	TimesPath   string
	Log         LogConfig
	Ftps        []FtpConfig
}

func processConfig(config fileConfig) *ExtractorConfig {
	var newConfig = &ExtractorConfig{
		Log2Console: config.Log2Console,
		TimesPath:   config.TimesPath,
		Log: LogConfig{
			Folder:   config.Log.Folder,
			Filename: config.Log.Filename,
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
		match := ServiceConfig{
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
	return newConfig
}

func ReadConfig(filepath string) (*ExtractorConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config fileConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}
	// replace logfile name data placeholder
	startIndex := strings.Index(config.Log.Filename, "{")
	endIndex := strings.Index(config.Log.Filename, "}")
	if startIndex != -1 && endIndex != -1 {
		dateFormat := config.Log.Filename[startIndex+1 : endIndex]
		config.Log.Filename = fmt.Sprintf(
			"%s%s%s",
			config.Log.Filename[:startIndex],
			time.Now().UTC().Format(dateFormat),
			config.Log.Filename[endIndex+1:],
		)
	}

	// process read config to match/duplicate service and server
	return processConfig(config), nil
}

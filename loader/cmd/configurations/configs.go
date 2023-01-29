package configurations

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ModeLogConfig struct {
	Folder   string `yaml:"folder"`
	Filename string `yaml:"filename"`
}
type LogConfig struct {
	Import ModeLogConfig `yaml:"compress"`
	Export ModeLogConfig `yaml:"cleaner"`
}

type ServiceConfig struct {
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	Src       string `yaml:"sourceFolder"`
	Dst       string `yaml:"destinationFolder"`
	Prefix    string `yaml:"filePrefix"`
	Extension string `yaml:"fileExtension"`
	MaxTime   int    `yaml:"maxTime"`
	Window    int    `yaml:"windowLimit"`
}

type loaderConfig struct {
	Log2Console bool            `yaml:"log2console"`
	TimesPath   string          `yaml:"timesFilePath"`
	Log         LogConfig       `yaml:"logging"`
	Services    []ServiceConfig `yaml:"services"`
}

func ReadConfig(filepath string) (*loaderConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var loader loaderConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &loader)
	if err != nil {
		return nil, err
	}

	return &loader, nil
}

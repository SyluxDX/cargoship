package configurations

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

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

type LoaderConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"timesFilePath"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Services []ServiceConfig `yaml:"services"`
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

func ReadConfig(filepath string) (*LoaderConfig, error) {
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

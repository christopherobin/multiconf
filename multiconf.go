package multiconf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	// 3rd parties
	"github.com/christopherobin/go-appdirs"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Parser func([]byte) (map[string]interface{}, error)

type Multiconf struct {
	AppConf appdirs.AppConf
	Files   []string
	Root    map[string]interface{}
	Parser  Parser
}

// Creates a new multiconf instance, the application name is expected to find out
// your os specific data and configuration folders, the parser takes in the config
// data and returns a map[string]interface{} that is then merged on top of other
// configurations.
//
// The file names are text/template templates, the following variables are available
// in the template:
//   {{.Home}}       -> The home directory for the current user
//   {{.Data}}       -> The data directory for your app (eg ~/.local/share/appname)
//   {{.SiteData}}   -> The global data directory for your app (eg /usr/local/share/appname)
//   {{.Config}}     -> The config directory for your app (eg ~/.config/appname)
//   {{.SiteConfig}} -> The global config directory for your app (eg /etc/xdg/appname)
//   {{.Cache}}      -> The cache directory for your app (eg ~/.cache/appname)
//   {{.Log}}        -> The data directory for your app (eg ~/.cache/appname/logs)
//
// See the documentation for github.com/christopherobin/go-appdirs for more details
// on the folders generated for your current distro
func NewMulticonf(appName string, parser Parser, files ...string) *Multiconf {
	return &Multiconf{
		AppConf: appdirs.AppConf{Name: appName},
		Files:   files,
		Root:    make(map[string]interface{}),
		Parser:  parser,
	}
}

// Load the configuration
func (conf *Multiconf) Load() error {
	directories, err := conf.AppConf.Directories()
	if err != nil {
		return err
	}

	// reset root
	conf.Root = make(map[string]interface{})

	for idx, config := range conf.Files {
		fileNameTmpl, err := template.New(fmt.Sprintf("multiconf.%d", idx)).Parse(config)
		if err != nil {
			return err
		}
		var b bytes.Buffer
		err = fileNameTmpl.Execute(&b, directories)
		if err != nil {
			return err
		}
		fileName := b.String()

		// if the file does not exist, ignore it
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			continue
		}

		rawConfig, err := ioutil.ReadFile(fileName)
		if err != nil {
			continue
		}

		localConf, err := conf.Parser(rawConfig)
		if err != nil {
			return err
		}

		_ = mergo.Merge(&conf.Root, localConf)
	}

	return nil
}

// A default YAML parser
func YamlParser(in []byte) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := yaml.Unmarshal(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// A default JSON parser
func JsonParser(in []byte) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := json.Unmarshal(in, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

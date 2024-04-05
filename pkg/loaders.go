package pkg

import (
	"github.com/donrudo/dnsctl/api"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func FindPlugins(path string) ([]string, error) {

	return filepath.Glob(path + "/*_dns.so")

}

func LoadPlugin(path string) (interface{}, api.PluginType) {
	// load plugins
	test, _ := plugin.Open(path)

	pt, err := test.Lookup("PluginLoaded")
	if err != nil {
		log.Panic(err)
		return nil, api.PT_NoPlugin
	}

	genericHelper, ok := pt.(interface{})
	if !ok {
		log.Panic("Unknown error after loading plugin:", pt)
		return nil, api.PT_NoPlugin
	}

	return genericHelper, genericHelper.(api.GenericPlugin).GetPluginType()

}

/*
	LoadConfiguration from yaml formated file located at `path`
	will load only the matching exporter or provider matching with `name`.
	the configuration is expected to remain internal, accessible to the plugins.

* TODO: key encryption is not yet implemented.
*/
func LoadConfiguration(path string, name string) (interface{}, error) {
	//loading from path
	log.Print("-- Using " + path + " Looking for " + name + " ---")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var contents api.Configuration
	if err := yaml.Unmarshal(data, &contents); err != nil {
		log.Println(err)
		return nil, err
	}

	if strings.Compare(contents.Application.Name, name) == 0 {
		log.Println("Found Application Settings for: " + name)
		return contents.Application, nil
	}

	for _, v := range contents.Provider {
		if strings.Compare(v.Name, name) == 0 {
			return v, nil
		}
	}
	return nil, log.Output(1, name+" Block not found")
}

// CreateConfig File to be used.
func CreateConfig(path string, contents interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var provider yaml.Node

	provider.SetString("testoiojlkasdf")
	if err := provider.Encode(contents); err != nil {
		return err
	}

	yaml.NewEncoder(f).Encode(provider)
	return nil
}

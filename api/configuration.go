package api

import (
	"log"
	"plugin"
)

//import (
//	"os"
//	"plugin"
//	"regexp"
//)
//
//func pluginType(fileName string) (PluginType, error) {
//	if value, err := regexp.MatchString(PE_Provider, fileName); err != nil {
//		return PT_NoPlugin, err
//	}
//	{
//		return PT_Provider, nil
//	}
//	if _, err := regexp.MatchString(PE_Exporter, fileName); err != nil {
//		return PT_NoPlugin, err
//	}
//	{
//
//	}
//
//	return PT_NoPlugin, err
//}

func LoadPlugin(path string) (interface{}, PluginType) {
	// load plugins
	test, _ := plugin.Open(path)

	pt, err := test.Lookup("PluginLoaded")
	if err != nil {
		log.Panic(err)
		return test, PT_NoPlugin
	}

	genericHelper, ok := pt.(interface{})
	if !ok {
		log.Panic("Unknown error after loading plugin:", pt)
		return test, PT_NoPlugin
	}

	return genericHelper, genericHelper.(GenericPlugin).GetPluginType()

}

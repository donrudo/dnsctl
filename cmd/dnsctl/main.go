package main

import (
	"context"
	"flag"
	"github.com/donrudo/dnsctl/api"
	"github.com/donrudo/dnsctl/pkg"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	//// "setup, remove or update a domain configuration."
	//flagSetDomain = flag.NewFlagSet("add", flag.ExitOnError)
	//flagSetRecord = flag.NewFlagSet("list", flag.ExitOnError)
	//flagSetUpdate = flag.NewFlagSet("update", flag.ExitOnError)
	//flagSetDelete = flag.NewFlagSet("delete", flag.ExitOnError)
	//
	//flagAddRecord  = flagSetRecord.String("add", "", "Adds a record to a previously configured domain.")
	//flagShowDomain = flagSetDomain.String("show", "", "Lists the full list of records for the given domain.")
	//flagShowRecord = flagSetRecord.String("show", "", "Resolves or prints the value for given FQDN")
	////	flagDeleteRecord  = flagSetRecord.String("delete", "", "deletes given record. Format requires a FQDN")
	////	flagConfigShow = flagSetConfig.Bool("show", false, "Shows the configuration contents without API credentials.")
	//

	config  = flag.String("config", "configs/dnsctl.yaml", "Path to config file")
	version = flag.Bool("version", false, "Prints the version number and exits.")
	debug   = flag.Bool("debug", false, "Enables debug using port 8080")
	//perr  = log.New(os.Stderr, "", 0)
)

var Version string

type SupportedPlugin struct {
	ctx       context.Context
	providers []api.ProviderPlugin
	exporters []api.ExporterPlugin
}

func main() {
	var feature SupportedPlugin
	defaultLocation := "configs/dnsctl.yaml"
	flag.Parse()
	//flagSetDomain.Parse(os.Args[2:])
	//flagSetRecord.Parse(os.Args[2:])
	////flagSetConfig.Parse(os.Args[2:])

	if *version {
		println("Version:", Version)
		os.Exit(0)
	}
	if *debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	}
	if strings.Compare(*config, "") == 0 {
		*config = defaultLocation
	}

	// load configuration from file.
	ifaceAppConfig, err := pkg.LoadConfiguration(*config, "dnsctl")
	if err != nil {
		log.Println("Error loading " + *config + ": ")
		log.Panicln("Could not load configuration file: " + err.Error())
	}
	log.Printf("Configuration for %s Loaded: ", ifaceAppConfig.(api.ApplicationCfg).Name)
	log.Printf("--Loading plugins at: %s", ifaceAppConfig.(api.ApplicationCfg).PluginDir)

	strAllPlugins, err := pkg.FindPlugins(ifaceAppConfig.(api.ApplicationCfg).PluginDir)
	if err != nil {
		log.Fatal("Plugins not found: " + err.Error() + ifaceAppConfig.(api.ApplicationCfg).PluginDir)
	}

	//Start loading plugins
	var genericPlugin interface{}
	var pt api.PluginType
	log.Println(strAllPlugins)
	for _, v := range strAllPlugins {
		log.Printf(" === looking for %s", v)
		genericPlugin, pt = pkg.LoadPlugin(v)
		switch pt {
		case api.PT_Exporter:

			log.Print("Loaded Exporter plugin")
			break
		case api.PT_Provider:

			feature.providers = append(feature.providers, genericPlugin.(api.ProviderPlugin))

			if err := feature.providers[0].Init(feature.ctx, *config); err != nil {
				log.Panic("Error while Initializing Provider: ", err)
			}

			log.Printf("Loaded %s-%s Provider plugin: %s", feature.providers[0].GetProvider().Name, feature.providers[0].GetVersion(), feature.providers[0].GetProvider().Domain)
			break
		case api.PT_NoPlugin:
			log.Fatal("Plugin is in not compatible")
		}
	}

	// 1.2 Use API to perform requested actions.
	// For now we are just listing records.
	//TODO: move into individual switch depending the used flags and query.
	records, err := feature.providers[0].ListRecords()
	if err != nil {
		log.Fatalf("Error loading records from DNS Provider")
	}
	for _, v := range records {
		log.Printf("%s \t %s \t %s ", v.Name, v.Type, v.Value)
	}

	// 2. Fallback to read-only actions If no Config file is present or domain is not found there. (common resolve)
	// 2.1 If action required Update or Creation of records return Error.
	// 2.2 Use normal domain procedures for read-only tasks.
	// 3 Return output according to export plugin used.

}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/donrudo/dnsctl/api"
	"github.com/donrudo/dnsctl/pkg"
	"log"
	"net/http"
	"os"
	"strings"
)

type SupportedPlugin struct {
	ctx       context.Context
	domain    map[string][]string
	providers map[string]api.ProviderPlugin
	output    map[string]api.OutputPlugin
}

var (
	//// "setup, remove or update a domain configuration."
	//flagSetDomain = flag.NewFlagSet("add", flag.ExitOnError)
	//flagSetRecord = flag.NewFlagSet("list", flag.ExitOnError)
	//flagSetUpdate = flag.NewFlagSet("update", flag.ExitOnError)
	//flagSetDelete = flag.NewFlagSet("delete", flag.ExitOnError)
	//
	//flagAddRecord  = flagSetRecord.String("add", "", "Adds a record to a previously configured domain.")
	//flagShowRecord = flagSetRecord.String("show", "", "Resolves or prints the value for given FQDN")
	////	flagDeleteRecord  = flagSetRecord.String("delete", "", "deletes given record. Format requires a FQDN")
	////	flagConfigShow = flagSetConfig.Bool("show", false, "Shows the configuration contents without API credentials.")
	//

	//verbose       = flag.Bool("verbose", false, "Enables verbosity for debug")
	listRecords   = flag.String("list", "", "List records for the given domain")
	config        = flag.String("config", "configs/dnsctl.yaml", "Path to config file")
	version       = flag.Bool("version", false, "Prints the version number and exits.")
	debug         = flag.Bool("debug", false, "Enables debug using port 8080")
	output        = flag.String("output", "Stdout", "Sets the Output plugin to be used")
	listOutputs   = flag.Bool("list-output", false, "Lists the Output plugin found")
	listProviders = flag.Bool("list-providers", false, "Lists the Provider plugins configured")
	listDomains   = flag.Bool("list-domains", false, "Lists the Domains loaded")
	//perr  = log.New(os.Stderr, "", 0)
)

var Version string
var feature SupportedPlugin

/*
ShowSettings Go over the flags used to display information.
*/
func ShowSettings() {

	if *listOutputs {
		var supportedOutputs []string
		for name, _ := range feature.output {
			supportedOutputs = append(supportedOutputs, name)
		}
		fmt.Println("Supported Output are:")
		fmt.Println("\t", strings.Join(supportedOutputs, " "))
	}
	if *listProviders {
		var supportedProviders []string
		for name, _ := range feature.providers {
			supportedProviders = append(supportedProviders, name)
		}
		fmt.Println("Configured Providers are:")
		fmt.Println("\t", strings.Join(supportedProviders, " "))
	}
	if *listDomains {
		fmt.Println("Configured Providers are:")
		var supportedDomains []string
		for name, _ := range feature.domain {
			supportedDomains = append(supportedDomains, name)
		}
		fmt.Println("Configured Providers are:")
		fmt.Println("\t", strings.Join(supportedDomains, " "))
	}
}

func main() {
	feature.domain = make(map[string][]string)
	feature.providers = make(map[string]api.ProviderPlugin)
	feature.output = make(map[string]api.OutputPlugin)

	defaultLocation := "configs/dnsctl.yaml"
	flag.Parse()
	//flagSetDomain.Parse(os.Args[2:])
	//flagSetRecord.Parse(os.Args[2:])
	////flagSetConfig.Parse(os.Args[2:])

	//if *verbose {
	//
	//}
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

	for _, v := range strAllPlugins {
		log.Printf(" === looking for %s", v)
		genericPlugin, pt = pkg.LoadPlugin(v)
		switch pt {
		case api.PtOutput:
			// Initializes generic plugins as OutputPlugin
			op, err := genericPlugin.(api.OutputPlugin).Init(feature.ctx, *config)
			if err != nil {
				log.Panic("Error while Initializing Provider: ", err)
			}

			feature.output[(op).GetName()] = op
			log.Print("Loaded Output plugin: ", op.GetName(), op.GetVersion())
			break
		case api.PtProvider:
			// Initializes generic plugins as ProviderPlugin.
			np, err := genericPlugin.(api.ProviderPlugin).Init(feature.ctx, *config)
			if err != nil {
				log.Println("Error while Initializing Provider: ", err)
			}

			feature.providers[(np).GetName()] = np
			log.Printf("Loaded %s-%s Provider plugin: %s",
				feature.providers[(np).GetName()].GetName(),
				feature.providers[(np).GetName()].GetVersion(),
				feature.providers[(np).GetName()].GetProvider().Domain)

			// Builds a helper Map created to prevent frequent scanning of the arrays
			//	while looking for a domain to be updated
			for _, dn := range (np).GetProvider().Domain {
				feature.domain[dn] = append(feature.domain[dn], (np).GetName())
			}

			break
		case api.PtNoPlugin:
			log.Fatal("Plugin is in not compatible")
		}
	}

	// Show Settings as requested by the user.
	ShowSettings()

	// 1.2 Use API to perform requested actions.
	// For now we are just listing records.
	//TODO: move into individual switch depending the used flags and query.

	if strings.Compare((*listRecords), "") > 0 {
		var recordValues []api.Record

		if arrProviders, ok := feature.domain[(*listRecords)]; ok {
			for _, strName := range arrProviders {
				recordValues, err = feature.providers[strName].ListRecordsFromDomain((*listRecords))
				if err != nil {
					log.Println("Error requesting records: ", err.Error())
				}
				log.Println("querying for record: ", recordValues)
				for _, record := range recordValues {
					if err := feature.output[(*output)].PrintRecord(record); err != nil {
						log.Fatal("Error at sending records to output: ", err.Error())
					}
				}
			}
		}
	}
	//records, err := feature.providers[0].ListRecords()
	//if err != nil {
	//	log.Fatalf("Error loading records from DNS Provider")
	//}
	//for _, v := range records {
	//	log.Printf("%s \t %s \t %s ", v.Name, v.Type, v.Value)
	//	if err := feature.output[0].PrintRecord(v); err != nil {
	//		log.Println(err.Error())
	//	}
	//}

	// 2. Fallback to read-only actions If no Config file is present or domain is not found there. (common resolve)
	// 2.1 If action required Update or Creation of records return Error.
	// 2.2 Use normal domain procedures for read-only tasks.
	// 3 Return output according to export plugin used.

}

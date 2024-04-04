package main

import (
	"context"
	"flag"
	"github.com/donrudo/dnsctl/api"
	"log"
	"net/http"
)

var (
	//// "setup, remove or update a domain configuration."
	//flagSetDomain = flag.NewFlagSet("add", flag.ExitOnError)
	//flagSetRecord = flag.NewFlagSet("list", flag.ExitOnError)
	//flagSetUpdate = flag.NewFlagSet("update", flag.ExitOnError)
	//flagSetDelete = flag.NewFlagSet("delete", flag.ExitOnError)
	//
	//flagAddRecord  = flagSetRecord.String("add", "", "Adds a a record to a previously configured domain.")
	//flagShowDomain = flagSetDomain.String("show", "", "Lists the full list of records for the given domain.")
	//flagShowRecord = flagSetRecord.String("show", "", "Resolves or prints the value for given FQDN")
	////	flagDeleteRecord  = flagSetRecord.String("delete", "", "deletes given record. Format requires a FQDN")
	////	flagConfigShow = flagSetConfig.Bool("show", false, "Shows the configuration contents without API credentials.")
	//
	debug = flag.Bool("debug", false, "Enables debug using port 8080")
	//perr  = log.New(os.Stderr, "", 0)
)

type SupportedPlugin struct {
	ctx       context.Context
	providers []api.ProviderPlugin
	exporters []api.ExporterPlugin
}

func main() {
	var feature SupportedPlugin

	// var err error
	flag.Parse()
	//flagSetDomain.Parse(os.Args[2:])
	//flagSetRecord.Parse(os.Args[2:])
	////flagSetConfig.Parse(os.Args[2:])

	if *debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	}

	genericPlugin, pt := api.LoadPlugin("build/plugins/cloudflare_dns.so")
	switch pt {
	case api.PT_Exporter:

		log.Print("Loaded Exporter plugin")
		break
	case api.PT_Provider:

		feature.providers = append(feature.providers, genericPlugin.(api.ProviderPlugin))

		if err := feature.providers[0].Init(feature.ctx); err != nil {
			log.Panic("Error while Initializing Provider: ", err)
		}

		log.Printf("Loaded %s Provider plugin: %s", feature.providers[0].GetProvider().Name, feature.providers[0].GetProvider().Domain)
		break
	case api.PT_NoPlugin:
		log.Fatal("Plugin is in not compatible")
	}

	// Process parameters.
	// Get Domain from FQDN.
	// 1. Load Configuration file.
	// 1.1 Search at config for API based on domain name.
	// 1.2 Use API to perform requested actions.

	// 2. Fallback to read-only actions If no Config file is present or domain is not found there.
	// 2.1 If action required Update or Creation of records return Error.
	// 2.2 Use normal domain procedures for read-only tasks.
	// 3 Return output according to export plugin used.

}

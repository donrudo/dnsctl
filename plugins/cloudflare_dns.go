package main

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/donrudo/dnsctl/api"
	"github.com/donrudo/dnsctl/pkg"
	"log"
)

var Version string
var PluginLoaded cloudflareProvider

type cloudflareProvider struct {
	api.ProviderPlugin
	ctx        context.Context
	connection *cloudflare.API
	Provider   api.SettingsProvider
}

func (t *cloudflareProvider) GetVersion() string {
	return Version
}

func (t *cloudflareProvider) GetPluginType() api.PluginType {
	return api.PtProvider
}

func (t *cloudflareProvider) GetProvider() api.SettingsProvider {
	return t.Provider
}

func (t *cloudflareProvider) GetName() string {
	return t.Provider.Name
}

/*
ListRecordsFromDomain from all the configured domains for Cloudflare.

	Returns: []api.Record with all the records found per domain.
	err is returned if the ID for a domain cannot be found or if records for at least one domain cannot be requested.
*/
func (t *cloudflareProvider) ListRecordsFromDomain(domain string) ([]api.Record, error) {
	// Fetch the zone ID
	var allRecords []api.Record
	id, err := t.connection.ZoneIDByName(domain)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cfRecords, resultInfo, err := t.connection.ListDNSRecords(t.ctx,
		cloudflare.ZoneIdentifier(id),
		cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	t.ctx = context.Background()
	log.Println(resultInfo)

	for _, cfr := range cfRecords {

		allRecords = append(allRecords,
			api.Record{
				Name:     cfr.Name,
				TTL:      cfr.TTL,
				Value:    cfr.Content,
				Type:     cfr.Type,
				Provider: t.GetProvider().Name,
			})

	}

	return allRecords, nil
}

/*
ListRecords from all the configured domains for Cloudflare.

	Returns: []api.Record with all the records found per domain.
	err is returned if the ID for a domain cannot be found or if records for at least one domain cannot be requested.
*/
func (t *cloudflareProvider) ListRecords() ([]api.Record, error) {
	// Fetch the zone ID
	var allRecords []api.Record
	for _, v := range t.Provider.Domain {
		id, err := t.connection.ZoneIDByName(v)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		cfRecords, resultInfo, err := t.connection.ListDNSRecords(t.ctx,
			cloudflare.ZoneIdentifier(id),
			cloudflare.ListDNSRecordsParams{})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		t.ctx = context.Background()
		log.Println(resultInfo)

		for _, cfr := range cfRecords {

			allRecords = append(allRecords,
				api.Record{
					Name:     cfr.Name,
					TTL:      cfr.TTL,
					Value:    cfr.Content,
					Type:     cfr.Type,
					Provider: t.GetProvider().Name,
				})
		}
	}
	return allRecords, nil
}
func (t *cloudflareProvider) Init(ctx context.Context, configFile string) (api.ProviderPlugin, error) {

	cloudflarePlugin, err := t.InitCloudflare(ctx, configFile)

	return &cloudflarePlugin, err

}

/*
	 cloudflareProvider Init, initializes plugin from configuration file
	    1. Load Configuration file.
	 	1.1 Search at config for API based on domain name.
		Stores basic settings for the applications, discards key and email after loading API
*/
func (t *cloudflareProvider) InitCloudflare(ctx context.Context, configFile string) (cloudflareProvider, error) {

	ProviderCfg, err := pkg.LoadConfiguration(configFile, "Cloudflare")
	if err != nil {
		log.Fatal(err)
		return *t, err
	}
	t.Provider = api.SettingsProvider{
		ProviderCfg.(api.ProviderCfg).Name,
		ProviderCfg.(api.ProviderCfg).Domain,
	}

	t.ctx = ctx
	if t.connection, err = cloudflare.New(ProviderCfg.(api.ProviderCfg).Key, ProviderCfg.(api.ProviderCfg).Email); err != nil {
		return *t, err
	}
	t.ctx = context.Background()

	// Most API calls require a Context
	return *t, nil
}

// Useful for test snippets
//func main() {
//	var feature cloudflareProvider
//	var ctx context.Context
//	if err := feature.Init(ctx); err != nil {
//		log.Print(err)
//	}
//}

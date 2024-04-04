package main

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/donrudo/dnsctl/api"
	"log"
	"os"
)

type cloudflareProvider struct {
	api.ProviderPlugin
	ctx        context.Context
	connection *cloudflare.API
	Provider   *api.Provider
}

func (t *cloudflareProvider) GetPluginType() api.PluginType {
	return api.PT_Provider
}
func (t *cloudflareProvider) GetProvider() *api.Provider {
	return t.Provider
}
func (t *cloudflareProvider) ListRecords() ([]api.Record, error) {
	// Fetch the zone ID
	id, err := t.connection.ZoneIDByName(t.Provider.Domain)
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

	var all_records []api.Record
	for _, cfr := range cfRecords {

		all_records = append(all_records,
			api.Record{
				cfr.Name,
				cfr.TTL,
				cfr.Content,
				cfr.Type,
				t.GetProvider().Name,
			})

	}

	return all_records, nil
}

func (t *cloudflareProvider) Init(ctx context.Context) error {

	//TODO: Load credentials, provider name and domain from configuration file before pushing
	t.Provider.Name = "Cloudflare"
	t.Provider.Domain = "example.com"
	t.ctx = ctx
	var err error
	if t.connection, err = cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL")); err != nil {
		return err
	}
	t.ctx = context.Background()
	// Most API calls require a Context
	return nil
}

var PluginLoaded cloudflareProvider

// Useful for test snippets
//func main() {
//	var feature cloudflareProvider
//	var ctx context.Context
//	if err := feature.Init(ctx); err != nil {
//		log.Print(err)
//	}
//}

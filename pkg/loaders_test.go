package pkg

import (
	"github.com/donrudo/dnsctl/api"
	"os"
	"testing"
)

func TestCreateConfig(t *testing.T) {
	tmpdir := t.TempDir()
	filename := "/dnsctl.yaml"

	Providers := make([]api.ProviderCfg, 0)
	Providers = append(Providers, api.ProviderCfg{
		Name:   "Providername",
		Domain: []string{"example.com.mx", "example.org", "example.com"},
		Key:    "keysldkjasdklaj",
		Email:  "user@email.com",
	})
	Providers = append(Providers, api.ProviderCfg{
		Name:   "Secondname",
		Domain: []string{"example.com.mx", "example.org", "example.com"},
		Key:    "sdfsdfsdfassdf",
		Email:  "user2@email.com",
	})
	ConfigFileSkell := api.Configuration{
		Provider: Providers}

	err := CreateConfig(tmpdir+filename, ConfigFileSkell)

	if err != nil {
		t.Fatalf(`CreateConfig("filepath") = %v, want "", error`, err)
	}

	if _, err := os.ReadFile(tmpdir + filename); err != nil {
		t.Fatalf(`CreateConfig("filepath") = %v, want "", error`, err)
	}

}

func TestLoadConfiguration(t *testing.T) {

	strPath := "../configs/dnsctl_example.yaml"
	strProviderName := "Providername"

	if _, err := os.ReadFile(strPath); err != nil {
		t.Skipf("Couldn't load yaml example file: %s, %v", strPath, err)
	}

	cfgYaml, err := LoadConfiguration(strPath, strProviderName)
	if err != nil || cfgYaml == nil {
		t.Fatalf(`LoadConfiguration("filepath", "provider") %q = %v, want "api.ProviderCfg, nil", error`, cfgYaml, err)
	}

}

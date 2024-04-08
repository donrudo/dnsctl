package api

import "context"

type PluginType int8

const (
	PtNoPlugin PluginType = iota
	PtOutput
	PtProvider
)

const (
	PeOutput   = "/*out.so"
	PeProvider = "/*dns.so"
)

type GenericPlugin interface {
	GetPluginType() PluginType
	GetVersion() string
	GetName() string
}

type ProviderPlugin interface {
	GenericPlugin
	ListRecords() ([]Record, error)
	ListRecordsFromDomain(domain string) ([]Record, error)
	GetProvider() SettingsProvider
	Init(context.Context, string) (ProviderPlugin, error)
}

type OutputPlugin interface {
	GenericPlugin
	PrintRecord(r Record) error
	PrintProvider(p SettingsProvider) error
	Init(context.Context, string) (OutputPlugin, error)
}

type Record struct {
	Name     string
	TTL      int
	Value    string
	Type     string
	Provider string
}

type SettingsProvider struct {
	Name   string   `yaml:"name"`
	Domain []string `yaml:"domain"`
}

type Configuration struct {
	Application ApplicationCfg `yaml:"application"`
	Provider    []ProviderCfg  `yaml:"provider"`
}

type ProviderCfg struct {
	Name   string   `yaml:"name"`
	Domain []string `yaml:"domain"`
	Key    string   `yaml:"credentials.key"`
	Email  string   `yaml:"credentials.email"`
}
type ApplicationCfg struct {
	Name      string `yaml:"name"`
	PluginDir string `yaml:"path.plugins"`
}

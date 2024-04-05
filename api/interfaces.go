package api

import "context"

type PluginType int8

const (
	PT_NoPlugin PluginType = iota
	PT_Exporter
	PT_Provider
)

const (
	PE_Exporter = ".*exp.so"
	PE_Provider = ".*dns.so"
)

type GenericPlugin interface {
	GetPluginType() PluginType
	GetVersion() string
}

type ProviderPlugin interface {
	GenericPlugin
	ListRecords() ([]Record, error)
	GetProvider() SettingsProvider
	Init(context.Context, string) error
}

type ExporterPlugin interface {
	GenericPlugin
	ListRecords() []Record
	GetProvider() *SettingsProvider
	Init(context.Context, string) error
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

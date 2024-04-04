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
}

type ProviderPlugin interface {
	GenericPlugin
	ListRecords() ([]Record, error)
	GetProvider() *Provider
	Init(context.Context) error
}
type ExporterPlugin interface {
	GenericPlugin
	ListRecords() []Record
	GetProvider() *Provider
	Init(context.Context) error
}

type Record struct {
	Name     string
	TTL      int
	Value    string
	Type     string
	Provider string
}

type Provider struct {
	Name   string
	Domain string
}

package main

import (
	"context"
	"fmt"
	"github.com/donrudo/dnsctl/api"
)

var Version string
var PluginLoaded stdOutput

const (
	Name = "Stdout"
)

type stdOutput struct {
	api.OutputPlugin
}

/*
PrintRecord using the Generic Record values, not provider specific values.
*/
func (t *stdOutput) PrintRecord(r api.Record) error {
	fmt.Println(r.Name, r.Type, r.Value, r.TTL, r.Provider)
	return nil
}

func (t *stdOutput) PrintProvider(p api.SettingsProvider) error {
	fmt.Println(p.Name, p.Domain)
	return nil
}

func (t *stdOutput) GetName() string {
	return Name
}

func (t *stdOutput) GetVersion() string {
	return Version
}

func (t *stdOutput) GetPluginType() api.PluginType {
	return api.PtOutput
}

func (t *stdOutput) Init(ctx context.Context, config string) (api.OutputPlugin, error) {
	result, err := t.InitStdout(ctx, config)
	return &result, err
}

func (t *stdOutput) InitStdout(ctx context.Context, config string) (stdOutput, error) {

	println("--- Using Output: ", Name)
	ctx = context.Background()

	return *t, nil
}

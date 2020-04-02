package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-provider-local/local"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: local.Provider})
}

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-local/local"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: local.Provider})
}

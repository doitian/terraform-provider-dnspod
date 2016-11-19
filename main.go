package main

import (
	"github.com/3pjgames/terraform-provider-dnspod/dnspod"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dnspod.Provider,
	})
}

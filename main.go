package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.forge.orange-labs.fr/vrabrmc/terraform-provider-vra7/vrealize"
	"gitlab.forge.orange-labs.fr/vrabrmc/terraform-provider-vra7/utils"
)

func main() {
	utils.InitLog()
	opts := plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return vrealize.Provider()
		},
	}

	plugin.Serve(&opts)
}

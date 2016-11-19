package dnspod

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"
)

func Provider() terraform.ResourceProvider {
	return ProviderWithConfig(nil)
}

func ProviderWithConfig(c *client.Config) terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"login_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSPOD_LOGIN_TOKEN", nil),
				Description: "DNSPod Login Token",
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSPOD_ENDPOINT", nil),
				Description: "DNSPod API Endpoint",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dnspod_domain": resourceDomain(),
			"dnspod_record": resourceRecord(),
		},

		ConfigureFunc: providerConfigure(c),
	}
}

func providerConfigure(c *client.Config) func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := client.Config{}
		if c != nil {
			config = *c
		}

		if config.HttpClient == nil {
			config.HttpClient = &http.Client{}
		}
		if config.LoginToken == "" {
			config.LoginToken = d.Get("login_token").(string)
		}
		if config.Endpoint == "" {
			config.Endpoint = d.Get("endpoint").(string)
		}

		return config.Client()
	}
}

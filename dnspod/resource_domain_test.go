package dnspod

import (
	"fmt"
	"testing"

	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceDomain(t *testing.T) {
	var domain client.Domain

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "dnspod_domain.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists("dnspod_domain.foo", &domain),
					resource.TestCheckResourceAttr("dnspod_domain.foo", "domain", "terraform-provider-dnspod.org"),
				),
			},
			// Repeat the config
			resource.TestStep{
				Config: testAccDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDomainExists("dnspod_domain.foo", &domain),
					resource.TestCheckResourceAttr("dnspod_domain.foo", "domain", "terraform-provider-dnspod.org"),
				),
			},
		},
	})
}

const testAccDomainConfig = `
resource "dnspod_domain" "foo" {
	domain = "terraform-provider-dnspod.org"
}
`

func testAccCheckDomainDestroy(s *terraform.State) error {
	return testAccCheckDomainDestroyWithProvider(s, testAccProvider)
}

func testAccCheckDomainDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	apiClient := provider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dnspod_domain" {
			continue
		}

		var resp client.DomainInfoResponse
		err := apiClient.Call("Domain.Info", &client.DomainInfoRequest{DomainId: rs.Primary.ID}, &resp)
		if err == nil {
			return fmt.Errorf("Found domain not removed: %+v", resp.Domain)
		} else {
			if bsce, ok := err.(*client.BadStatusCodeError); !ok || bsce.Code != "6" {
				return err
			}
		}
	}

	return nil
}

func testAccCheckDomainExists(n string, i *client.Domain) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckDomainExistsWithProviders(n, i, &providers)
}

func testAccCheckDomainExistsWithProviders(n string, i *client.Domain, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			apiClient := provider.Meta().(*client.Client)
			var resp client.DomainInfoResponse
			err := apiClient.Call("Domain.Info", &client.DomainInfoRequest{DomainId: rs.Primary.ID}, &resp)
			if err != nil {
				return err
			}

			*i = resp.Domain
			return nil
		}

		return fmt.Errorf("Domain not found")
	}
}

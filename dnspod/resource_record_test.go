package dnspod

import (
	"fmt"
	"strings"
	"testing"

	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "dnspod_record.www",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccRecordConfig, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("dnspod_record.www"),
					resource.TestCheckResourceAttr("dnspod_domain.foo", "domain", testDomain),
					resource.TestCheckResourceAttr("dnspod_record.www", "sub_domain", "www"),
					resource.TestCheckResourceAttr("dnspod_record.www", "record_type", "A"),
					resource.TestCheckResourceAttr("dnspod_record.www", "value", "8.8.8.8"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccRecordConfig, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("dnspod_record.www"),
					resource.TestCheckResourceAttr("dnspod_domain.foo", "domain", testDomain),
					resource.TestCheckResourceAttr("dnspod_record.www", "sub_domain", "www"),
					resource.TestCheckResourceAttr("dnspod_record.www", "record_type", "A"),
					resource.TestCheckResourceAttr("dnspod_record.www", "value", "8.8.8.8"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccRecordConfigModify, testDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists("dnspod_record.www"),
					resource.TestCheckResourceAttr("dnspod_domain.foo", "domain", testDomain),
					resource.TestCheckResourceAttr("dnspod_record.www", "sub_domain", "www"),
					resource.TestCheckResourceAttr("dnspod_record.www", "record_type", "A"),
					resource.TestCheckResourceAttr("dnspod_record.www", "value", "8.8.4.4"),
				),
			},
		},
	})
}

const testAccRecordConfig = `
resource "dnspod_domain" "foo" {
	domain = "%s"
}

resource "dnspod_record" "www" {
	domain_id = "${dnspod_domain.foo.id}"
	sub_domain = "www"
	record_type = "A"
	value = "8.8.8.8"
}
`

const testAccRecordConfigModify = `
resource "dnspod_domain" "foo" {
	domain = "%s"
}

resource "dnspod_record" "www" {
	domain_id = "${dnspod_domain.foo.id}"
	sub_domain = "www"
	record_type = "A"
	value = "8.8.4.4"
}
`

func testAccCheckRecordExists(n string) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckRecordExistsWithProviders(n, &providers)
}

func testAccCheckRecordExistsWithProviders(n string, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		domainIdRecordId := strings.SplitN(rs.Primary.ID, "-", 2)
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			apiClient := provider.Meta().(*client.Client)
			var resp client.RecordInfoResponse
			err := apiClient.Call("Record.Info", &client.RecordInfoRequest{DomainId: domainIdRecordId[0], RecordId: domainIdRecordId[1]}, &resp)
			if err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("Record not found")
	}
}

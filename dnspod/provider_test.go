package dnspod

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testDomain string

func init() {
	config := &client.Config{Lang: "en"}
	if os.Getenv("DNSPOD_DEBUG") != "" {
		config.Logger = log.New(os.Stdout, "DNSPOD API: ", log.LstdFlags)
	}

	testAccProvider = ProviderWithConfig(config).(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"dnspod": testAccProvider,
	}
	testDomain = fmt.Sprintf("terraform-provider-dnspod-%d.org", time.Now().Unix())
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	loginToken := os.Getenv("DNSPOD_LOGIN_TOKEN")

	if loginToken == "" {
		t.Fatal("DNSPOD_LOGIN_TOKEN is not set")
	}
}

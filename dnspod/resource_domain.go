package dnspod

import (
	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	req := client.DomainCreateRequest{Domain: d.Get("domain").(string)}
	var resp client.DomainCreateResponse
	err := conn.Call("Domain.Create", &req, &resp)
	if err != nil {
		return err
	}

	id := resp.Domain.Id
	d.SetId(id)

	return nil
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	req := client.DomainInfoRequest{DomainId: d.Id()}
	var resp client.DomainInfoResponse
	err := conn.Call("Domain.Info", &req, &resp)
	if err != nil {
		if bsce, ok := err.(*client.BadStatusCodeError); ok && bsce.Code == "6" {
			// 6 域名ID错误
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("domain", resp.Domain.Name)

	return nil
}

func resourceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	var resp client.DomainRemoveResponse
	req := client.DomainRemoveRequest{DomainId: d.Id()}
	err := conn.Call("Domain.Remove", &req, &resp)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

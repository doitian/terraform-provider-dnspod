package dnspod

import (
	"strings"

	"github.com/3pjgames/terraform-provider-dnspod/dnspod/client"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sub_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_line": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "默认",
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mx": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  600,
			},
			"weight": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	req := client.RecordCreateRequest{}
	req.DomainId = d.Get("domain_id").(string)
	req.SubDomain = d.Get("sub_domain").(string)
	req.RecordType = d.Get("record_type").(string)
	req.RecordLine = d.Get("record_line").(string)
	req.Value = d.Get("value").(string)
	req.Mx = d.Get("mx").(string)
	req.Ttl = d.Get("ttl").(string)
	req.Weight = d.Get("weight").(string)

	var resp client.RecordCreateResponse
	err := conn.Call("Record.Create", &req, &resp)
	if err != nil {
		return err
	}

	id := (*resp.Record.Id).(string)
	d.SetId(req.DomainId + "-" + id)

	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	req := client.RecordModifyRequest{}
	req.DomainId, req.RecordId = splitId(d.Id())
	req.SubDomain = d.Get("sub_domain").(string)
	req.RecordType = d.Get("record_type").(string)
	req.RecordLine = d.Get("record_line").(string)
	req.Value = d.Get("value").(string)
	req.Mx = d.Get("mx").(string)
	req.Ttl = d.Get("ttl").(string)
	req.Weight = d.Get("weight").(string)

	var resp client.RecordModifyResponse
	err := conn.Call("Record.Modify", &req, &resp)
	if err != nil {
		if bsce, ok := err.(*client.BadStatusCodeError); ok && (bsce.Code == "6" || bsce.Code == "8") {
			// 6 域名ID错误
			// 8 记录ID错误
			d.SetId("")
			return nil
		}

		return err
	}

	return nil
}

func resourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	domainId, recordId := splitId(d.Id())
	req := client.RecordInfoRequest{RecordId: recordId, DomainId: domainId}
	var resp client.RecordInfoResponse
	err := conn.Call("Record.Info", &req, &resp)
	if err != nil {
		if bsce, ok := err.(*client.BadStatusCodeError); ok && (bsce.Code == "6" || bsce.Code == "8") {
			// 6 域名ID错误
			// 8 记录ID错误
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("domain_id", domainId)
	d.Set("sub_domain", resp.Record.SubDomain)
	d.Set("record_type", resp.Record.RecordType)
	d.Set("record_line", resp.Record.RecordLine)
	d.Set("value", resp.Record.Value)
	d.Set("mx", resp.Record.Mx)
	d.Set("ttl", resp.Record.Ttl)
	d.Set("weight", resp.Record.Weight)

	return nil
}

func resourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	var resp client.RecordRemoveResponse
	req := client.RecordRemoveRequest{DomainId: d.Id()}
	err := conn.Call("Record.Remove", &req, &resp)
	if err != nil {
		if bsce, ok := err.(*client.BadStatusCodeError); !ok || (bsce.Code != "6" && bsce.Code != "8") {
			return err
		}
	}

	d.SetId("")

	return nil
}

func splitId(id string) (string, string) {
	parts := strings.SplitN(id, "-", 2)
	return parts[0], parts[1]
}

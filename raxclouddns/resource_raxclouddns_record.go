package raxclouddns

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gophercloud/gophercloud/v2"

	"github.com/rackerlabs/goclouddns"
	"github.com/rackerlabs/goclouddns/records"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS domain where this record will be created",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name you want to create",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record type",
			},
			"data": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record data",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time-To-Live (TTL) for the record",
				Default:     0,
			},
			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "DNS record priority (MX and SRV only)",
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"comment": &schema.Schema{
				Type:        schema.TypeString,
				Description: "User specified comment",
				Optional:    true,
				Default:     "",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("domain_id", func(old, new, meta interface{}) bool {
				// if "domain_id" changes then we must delete the old resource and
				// create a new resource
				return new != old
			}),
		),
	}
}

func resourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	domID := d.Get("domain_id").(string)

	opts := records.CreateOpts{
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Data:    d.Get("data").(string),
		TTL:     uint(d.Get("ttl").(int)),
		Comment: d.Get("comment").(string),
	}

	if opts.Type == "MX" || opts.Type == "SRV" {
		opts.Priority = uint(d.Get("priority").(int))
	}

	record, err := records.Create(context.Background(), service, domID, opts).Extract()
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", domID, record.ID))
	return resourceRecordRead(d, m)
}

func resourceRecordRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	resID := strings.Split(d.Id(), ":")
	if len(resID) < 2 {
		return fmt.Errorf("Invalid resource ID '%s'", d.Id())
	}

	domID := resID[0]

	record, err := records.Get(context.Background(), service, domID, resID[1]).Extract()
	if err != nil {
		return err
	}

	d.Set("domain_id", domID)
	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("data", record.Data)
	d.Set("ttl", record.TTL)
	d.Set("priority", record.Priority)
	d.Set("comment", record.Comment)
	d.SetId(fmt.Sprintf("%s:%s", domID, record.ID))
	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	resID := strings.Split(d.Id(), ":")
	if len(resID) < 2 {
		return fmt.Errorf("Invalid resource ID '%s'", d.Id())
	}

	domID := resID[0]

	record, err := records.Get(context.Background(), service, domID, resID[1]).Extract()
	if err != nil {
		return err
	}

	opts := records.UpdateOpts{
		Name:     d.Get("name").(string),
		Data:     d.Get("data").(string),
		TTL:      uint(d.Get("ttl").(int)),
		Priority: uint(d.Get("priority").(int)),
		Comment:  d.Get("comment").(string),
	}

	err = records.Update(context.Background(), service, domID, record, opts).ExtractErr()
	if err != nil {
		return err
	}

	return resourceRecordRead(d, m)
}

func resourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	resID := strings.Split(d.Id(), ":")
	if len(resID) < 2 {
		return fmt.Errorf("Invalid resource ID '%s'", d.Id())
	}

	domID := resID[0]

	return records.Delete(context.Background(), service, domID, resID[1]).ExtractErr()
}

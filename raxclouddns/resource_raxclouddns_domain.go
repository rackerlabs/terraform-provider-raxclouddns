package raxclouddns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud/v2"

	"github.com/rackerlabs/goclouddns"
	"github.com/rackerlabs/goclouddns/domains"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name you want to create",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email address",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time-To-Live (TTL) for the domain",
				Default:     3600,
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
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	opts := domains.CreateOpts{
		Name:    d.Get("name").(string),
		Email:   d.Get("email").(string),
		TTL:     uint(d.Get("ttl").(int)),
		Comment: d.Get("comment").(string),
	}

	domain, err := domains.Create(context.Background(), service, opts).Extract()
	if err != nil {
		return err
	}

	d.SetId(domain.ID)
	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	domain, err := domains.Get(context.Background(), service, d.Id()).Extract()
	if err != nil {
		return err
	}

	d.Set("name", domain.Name)
	d.Set("ttl", domain.TTL)
	d.Set("account_id", domain.AccountID)
	d.Set("email", domain.EmailAddress)
	d.Set("created", domain.Created)
	d.Set("updated", domain.Updated)
	d.Set("comment", domain.Comment)
	d.SetId(domain.ID)
	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	domain, err := domains.Get(context.Background(), service, d.Id()).Extract()
	if err != nil {
		return err
	}

	opts := domains.UpdateOpts{
		Email:   d.Get("email").(string),
		TTL:     uint(d.Get("ttl").(int)),
		Comment: d.Get("comment").(string),
	}

	err = domains.Update(context.Background(), service, domain, opts).ExtractErr()
	if err != nil {
		return err
	}

	return resourceDomainRead(d, m)
}

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	return domains.Delete(context.Background(), service, d.Id()).ExtractErr()
}

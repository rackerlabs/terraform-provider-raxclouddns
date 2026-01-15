package raxclouddns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud/v2"

	"github.com/rackerlabs/goclouddns"
	"github.com/rackerlabs/goclouddns/domains"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCredRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain you want to access",
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCredRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	service, err := goclouddns.NewCloudDNS(config.OsClient, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	domName := d.Get("name").(string)
	if domName == "" {
		return fmt.Errorf("Cannot supply an empty domain name")
	}

	opts := domains.ListOpts{
		Name: domName,
	}

	allPages, err := domains.List(context.Background(), service, &opts).AllPages(context.Background())
	if err != nil {
		return err
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		return err
	}

	if len(allDomains) == 0 {
		return fmt.Errorf("No results found for '%s'", domName)
	}

	if len(allDomains) > 1 {
		return fmt.Errorf("Too many results found for '%s'", domName)
	}

	domain, err := domains.Get(context.Background(), service, allDomains[0].ID).Extract()
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

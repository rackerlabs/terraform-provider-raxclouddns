package raxclouddns

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: "The Identity authentication URL.",
				Default:     "https://identity.api.rackspacecloud.com/v2.0/",
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USERNAME",
				}, ""),
				Description: "Username to login with.",
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_PASSWORD",
				}, ""),
				Description: "Password to login with.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"raxclouddns_domain": resourceDomain(),
			"raxclouddns_record": resourceRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"raxclouddns_domain": dataSourceDomain(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		IdentityEndpoint: d.Get("auth_url").(string),
		Password:         d.Get("password").(string),
		Username:         d.Get("user_name").(string),
		TerraformVersion: terraformVersion,
		SDKVersion:       meta.SDKVersionString(),
	}

	if err := config.authenticate(); err != nil {
		return nil, err
	}

	log.Println("[INFO] Initialized raxclouddns client")
	return &config, nil
}

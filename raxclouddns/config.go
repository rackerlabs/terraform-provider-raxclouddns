package raxclouddns

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	tokens2 "github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	"github.com/rackerlabs/goraxauth"
)

type Config struct {
	IdentityEndpoint string
	Username         string
	Password         string
	ApiKey           string

	OsClient *gophercloud.ProviderClient

	TerraformVersion string
	SDKVersion       string
}

func (c *Config) authenticate() error {
	if c.IdentityEndpoint == "" {
		return fmt.Errorf("'auth_url' must be left as the default or specified")
	}

	if c.Username == "" {
		return fmt.Errorf("No username supplied via 'user_name' or OS_USERNAME env")
	}

	if c.Password == "" && c.ApiKey == "" {
		return fmt.Errorf("Either 'password' (OS_PASSWORD) or 'api_key' (RAX_API_KEY) must be provided")
	}

	opts := goraxauth.AuthOptions{
		AuthOptions: tokens2.AuthOptions{
			IdentityEndpoint: c.IdentityEndpoint,
			Username:         c.Username,
			Password:         c.Password,
		},
		ApiKey: c.ApiKey,
	}

	client, err := goraxauth.AuthenticatedClient(context.Background(), opts)
	if err != nil {
		return err
	}

	ua := fmt.Sprintf("HashiCorp Terraform/%s Terraform Plugin SDK/%s", c.TerraformVersion, c.SDKVersion)
	client.UserAgent.Prepend(ua)

	c.OsClient = client
	return nil
}

package raxclouddns

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

type Config struct {
	IdentityEndpoint string
	Username         string
	Password         string

	OsClient *gophercloud.ProviderClient
	authOpts *gophercloud.AuthOptions

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

	if c.Password == "" {
		return fmt.Errorf("No password supplied via 'password' or OS_PASSWORD env")
	}

	authOpts := &gophercloud.AuthOptions{
		IdentityEndpoint: c.IdentityEndpoint,
		Username:         c.Username,
		Password:         c.Password,
	}

	client, err := openstack.NewClient(authOpts.IdentityEndpoint)
	if err != nil {
		return err
	}

	ua := fmt.Sprintf("HashiCorp Terraform/%s Terraform Plugin SDK/%s", c.TerraformVersion, c.SDKVersion)

	client.UserAgent.Prepend(ua)

	if err := openstack.Authenticate(context.Background(), client, *authOpts); err != nil {
		return err
	}

	c.OsClient = client
	c.authOpts = authOpts

	return nil
}

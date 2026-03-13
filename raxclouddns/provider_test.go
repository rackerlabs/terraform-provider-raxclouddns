package raxclouddns

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"raxclouddns": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestAccProviderConfigureWithApiKey verifies that the provider can
// authenticate using an API key instead of a password.
//
// Run with:
//
//	TF_ACC=1 go test -v -run TestAccProviderConfigureWithApiKey ./raxclouddns/
//
// Required env vars: OS_USERNAME, RAX_API_KEY
// TestAccProviderConfigure verifies that the provider can authenticate
// using either API key or password credentials.
//
// Run with:
//
//	TF_ACC=1 go test -v -run TestAccProviderConfigure ./raxclouddns/
//
// Required env vars: OS_USERNAME and one of RAX_API_KEY or OS_PASSWORD
func TestAccProviderConfigure(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping acceptance test")
	}

	username := os.Getenv("OS_USERNAME")
	if username == "" {
		t.Fatal("OS_USERNAME must be set for acceptance tests")
	}

	apiKey := os.Getenv("RAX_API_KEY")
	password := os.Getenv("OS_PASSWORD")

	if apiKey == "" && password == "" {
		t.Fatal("one of RAX_API_KEY or OS_PASSWORD must be set for acceptance tests")
	}

	raw := map[string]interface{}{
		"auth_url":  "https://identity.api.rackspacecloud.com/v2.0/",
		"user_name": username,
		"api_key":   apiKey,
		"password":  password,
	}

	authMethod := "password"
	if apiKey != "" {
		authMethod = "API key"
	}

	rawConfig := terraform.NewResourceConfigRaw(raw)

	err := testAccProvider.Configure(rawConfig)
	if err != nil {
		t.Fatalf("provider configure error (%s auth): %s", authMethod, err)
	}

	config, ok := testAccProvider.Meta().(*Config)
	if !ok {
		t.Fatal("expected provider meta to be *Config")
	}

	if config.OsClient == nil {
		t.Fatal("expected authenticated OsClient, got nil")
	}

	t.Logf("Successfully authenticated with %s", authMethod)
}

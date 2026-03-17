---
page_title: "Provider: Rackspace Cloud DNS"
description: |-
  The Rackspace Cloud DNS provider is used to manage DNS domains and records on Rackspace Cloud DNS.
---

# Rackspace Cloud DNS Provider

The Rackspace Cloud DNS provider is used to manage DNS domains and records on Rackspace Cloud DNS.

## Example Usage

```terraform
provider "raxclouddns" {
  user_name = "my-username"
  password  = "my-password"
}

resource "raxclouddns_domain" "example" {
  name  = "example.com"
  email = "admin@example.com"
  ttl   = 3600
}

resource "raxclouddns_record" "www" {
  domain_id = raxclouddns_domain.example.id
  name      = "www.example.com"
  type      = "A"
  data      = "203.0.113.10"
  ttl       = 300
}
```

## Authentication

The provider authenticates against the Rackspace Identity API. Credentials can be provided via provider configuration or environment variables.

### Environment Variables

- `OS_USERNAME` - Rackspace username
- `OS_PASSWORD` - Rackspace password
- `RAX_API_KEY` - Rackspace API key (alternative to password)
- `OS_AUTH_URL` - Identity endpoint (defaults to `https://identity.api.rackspacecloud.com/v2.0/`)

## Argument Reference

- `user_name` - (Required) Rackspace username. Can also be set with `OS_USERNAME`.
- `password` - (Optional, Sensitive) Rackspace password. Can also be set with `OS_PASSWORD`.
- `api_key` - (Optional, Sensitive) Rackspace API key. Can also be set with `RAX_API_KEY`.
- `auth_url` - (Optional) Identity authentication URL. Defaults to `https://identity.api.rackspacecloud.com/v2.0/`.

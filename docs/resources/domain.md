---
page_title: "raxclouddns_domain Resource - Rackspace Cloud DNS"
subcategory: ""
description: |-
  Manages a DNS domain in Rackspace Cloud DNS.
---

# raxclouddns_domain (Resource)

Manages a DNS domain in Rackspace Cloud DNS.

## Example Usage

```terraform
resource "raxclouddns_domain" "example" {
  name    = "example.com"
  email   = "admin@example.com"
  ttl     = 3600
  comment = "Managed by Terraform"
}
```

## Argument Reference

- `name` - (Required) The domain name to create.
- `email` - (Required) The email address associated with the domain.
- `ttl` - (Optional) Time-To-Live in seconds. Defaults to `3600`.
- `comment` - (Optional) A comment for the domain.

## Import

Domains can be imported using their ID:

```shell
terraform import raxclouddns_domain.example DOMAIN_ID
```

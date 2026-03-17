---
page_title: "raxclouddns_domain Data Source - Rackspace Cloud DNS"
subcategory: ""
description: |-
  Look up an existing DNS domain in Rackspace Cloud DNS by name.
---

# raxclouddns_domain (Data Source)

Look up an existing DNS domain in Rackspace Cloud DNS by name.

## Example Usage

```terraform
data "raxclouddns_domain" "example" {
  name = "example.com"
}

output "domain_id" {
  value = data.raxclouddns_domain.example.id
}
```

## Argument Reference

- `name` - (Required) The domain name to look up.

## Attribute Reference

- `id` - The domain ID.
- `ttl` - The domain TTL in seconds.
- `account_id` - The Rackspace account ID that owns the domain.
- `email` - The email address associated with the domain.
- `created` - The creation timestamp.
- `updated` - The last updated timestamp.
- `comment` - The domain comment.

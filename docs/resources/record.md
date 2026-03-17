---
page_title: "raxclouddns_record Resource - Rackspace Cloud DNS"
subcategory: ""
description: |-
  Manages a DNS record in Rackspace Cloud DNS.
---

# raxclouddns_record (Resource)

Manages a DNS record in Rackspace Cloud DNS.

## Example Usage

```terraform
resource "raxclouddns_record" "www" {
  domain_id = raxclouddns_domain.example.id
  name      = "www.example.com"
  type      = "A"
  data      = "203.0.113.10"
  ttl       = 300
}

resource "raxclouddns_record" "mail" {
  domain_id = raxclouddns_domain.example.id
  name      = "example.com"
  type      = "MX"
  data      = "mail.example.com"
  priority  = 10
}
```

## Argument Reference

- `domain_id` - (Required) The ID of the domain this record belongs to. Changing this forces a new resource.
- `name` - (Required) The DNS record name.
- `type` - (Required) The DNS record type (e.g. A, AAAA, CNAME, MX, SRV, TXT).
- `data` - (Required) The record data (e.g. IP address, hostname).
- `ttl` - (Optional) Time-To-Live in seconds. Defaults to `0` (uses domain TTL).
- `priority` - (Optional) Record priority, valid for MX and SRV records only. Must be between 0 and 65535. Defaults to `0`.
- `comment` - (Optional) A comment for the record.

## Import

Records can be imported using the composite ID format `DOMAIN_ID:RECORD_ID`:

```shell
terraform import raxclouddns_record.example DOMAIN_ID:RECORD_ID
```

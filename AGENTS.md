# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Build Commands

```bash
# Build the provider
go build -o terraform-provider-raxclouddns

# Run tests
go test ./...

# Vet and format check
go vet ./...
gofmt -l .
```

## CI/CD

**CI** (`.github/workflows/ci.yml`): Runs on push to main and PRs. Executes build, test, vet, and format checks.

**Releases** (`.github/workflows/release.yml`): Triggered by tags matching `v*`. Uses GoReleaser to build releases with checksums.

```bash
# Validate GoReleaser config
goreleaser check

# Test release locally (no publish)
goreleaser release --snapshot --clean
```

## Architecture

This is a Terraform provider for Rackspace Cloud DNS, built using the Terraform Plugin SDK v1.

### Key Dependencies
- `github.com/hashicorp/terraform-plugin-sdk` - Terraform plugin framework (v1 SDK)
- `github.com/gophercloud/gophercloud/v2` - OpenStack client library for authentication
- `github.com/rackerlabs/goclouddns` - Rackspace Cloud DNS API client

### Code Structure

**Entry Point**: `main.go` - Starts the plugin server and registers the provider

**Provider Package** (`raxclouddns/`):
- `provider.go` - Defines provider schema (auth_url, user_name, password) and registers resources/data sources
- `config.go` - Handles authentication against Rackspace Identity API using gophercloud
- `resource_raxclouddns_domain.go` - CRUD operations for DNS domains
- `resource_raxclouddns_record.go` - CRUD operations for DNS records (A, MX, SRV, etc.)
- `data_source_raxclouddns_domain.go` - Read-only lookup of existing domains by name

### Resource ID Format
- Domain IDs: Simple string ID from API
- Record IDs: Composite format `{domain_id}:{record_id}` - parsed via `strings.Split(d.Id(), ":")`

### Authentication
Provider authenticates via environment variables `OS_USERNAME` and `OS_PASSWORD`, or explicit provider config. Uses Rackspace Identity endpoint (default: `https://identity.api.rackspacecloud.com/v2.0/`).

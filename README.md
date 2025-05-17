# Terraform Provider for LiteLLM

This provider enables Terraform to manage LiteLLM resources including users, API keys, models, teams, and organizations.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using `go build -o terraform-provider-litellm`

## Using the provider

```hcl
provider "litellm" {
  api_key = var.litellm_api_key
}

# Example resources will be added here
```

## Developing the Provider

### Requirements

- Go >= 1.19

### Building

```shell
go build -o terraform-provider-litellm
```

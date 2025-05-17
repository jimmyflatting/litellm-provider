# 🚀 LiteLLM Terraform Provider Plan

## 1. Overview

We are building a Terraform provider for LiteLLM to enable users to manage LiteLLM resources declaratively. This includes creating, reading, updating, and deleting entities such as users, keys, models, teams, and organizations via Terraform.

The provider will be written in Go using the Terraform Plugin SDK v2, and it will use the OpenAPI specification (`./openapi.json`) to generate the API client.

## 2. Supported Resources

### Resources
- `litellm_user` – Represents a user in LiteLLM
- `litellm_key` – Represents an API key in LiteLLM
- `litellm_model` – Represents a model configuration in LiteLLM
- `litellm_team` – Represents a team
- `litellm_organization` – Represents an organization

### Data Sources
- `litellm_user` – Read-only user data
- `litellm_key` – Read-only key data
- `litellm_model` – Read-only model config
- `litellm_team` – Read-only team info
- `litellm_organization` – Read-only organization info

## 3. Authentication

Authentication will be done via API Key:

```hcl
provider "litellm" {
  api_key  = var.litellm_api_key
  endpoint = var.litellm_endpoint # Optional, defaults to https://api.litellm.io
}
```

### Supported configuration sources:
- Terraform provider block
- Environment variable `LITELLM_API_KEY`
- Default endpoint fallback

## 4. Project Structure & Setup

Create a new Go module: `go mod init github.com/jimmyflatting/terraform-provider-litellm`

### Structure:
```
/cmd
/internal
  /client     # Generated OpenAPI client
  /provider   # Provider logic
  /resources  # Individual resource logic
  /datasource # Individual data source logic
/examples
/docs
main.go
provider.go
```

## 5. OpenAPI Client Integration

- Use oapi-codegen or openapi-generator to generate a Go client
- Custom wrapper may be used to simplify handling auth and pagination
- Keep OpenAPI spec updated and re-generate client as part of CI/CD

## 6. Resource Implementation (CRUD)

### Each resource must support:
- Create: POST to create resource
- Read: GET to sync state with API
- Update: PUT/PATCH to modify resource
- Delete: DELETE to remove resource

### Each resource should:
- Store ID in Terraform state
- Validate required/optional attributes
- Handle API error responses gracefully

## 7. Data Source Implementation

Each data source will:
- Accept identifying parameters (e.g., name or ID)
- Call GET endpoint
- Populate read-only attributes
- Be used for reference-only scenarios

## 8. Schema Design

Use `terraform-plugin-sdk/v2/helper/schema`:
```go
"email": {
    Type:        schema.TypeString,
    Required:    true,
    Description: "The email address of the user.",
}
```

## 9. Provider Configuration Schema

```go
"api_key": {
    Type:        schema.TypeString,
    Required:    true,
    DefaultFunc: schema.EnvDefaultFunc("LITELLM_API_KEY", nil),
    Description: "API Key for authenticating with LiteLLM.",
},
"endpoint": {
    Type:        schema.TypeString,
    Optional:    true,
    Default:     "https://api.litellm.io",
    Description: "Base URL for the LiteLLM API.",
}
```

## 10. Error Handling Strategy

- Handle all 4xx/5xx responses
- Implement retry logic for 429 (rate limits) and transient 5xx errors
- Provide meaningful diagnostics to Terraform logs

## 11. State Management Strategy

- Store resource ID and relevant attributes in state
- On Read, sync with API and update state
- If resource no longer exists, return `d.SetId("")` to signal drift

## 12. Testing Strategy

- Unit tests for schema and logic
- Acceptance tests using terraform-plugin-testing
- Optional: set up mock server or test workspace against staging API
- Include tests for each resource's lifecycle and edge cases

## 13. Documentation

- Auto-generate documentation with tfplugindocs
- Write usage examples for each resource and data source
- Add to `/docs/resources/` and `/docs/data-sources/`

## 14. CI/CD

Use GitHub Actions to:
- Lint (golangci-lint)
- Build/test
- Run acceptance tests
- Auto-generate documentation
- Tag and publish releases

## 15. Release & Versioning

- Use SemVer
- Tag releases as v0.x.0 until stable
- Publish to the Terraform Registry
- Automatically generate GitHub releases and CHANGELOG from tags

## 16. Maintenance Strategy

### Repository should include:
- `README.md` – overview and examples
- `CHANGELOG.md` – track changes
- `CONTRIBUTING.md` – guide contributors
- `CODEOWNERS` – if applicable

Additional considerations:
- Define policy for handling issues and PRs
- Periodically update OpenAPI client as the LiteLLM API evolves

## 17. Steps to Build

1. Initialize Go module and repo structure
2. Generate client from openapi.json
3. Implement provider schema
4. Implement resource CRUD for each type
5. Implement data sources
6. Add unit + acceptance tests
7. Generate and write documentation
8. Set up GitHub Actions for CI/CD
9. Add README, CHANGELOG, and CONTRIBUTING files
10. Publish to Terraform Registry
11. Update docs and maintain releases
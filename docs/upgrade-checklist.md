# Upgrade Checklist

- `docker-compose.yml`
  - https://hub.docker.com/r/amazon/dynamodb-local
- `docker/go/Dockerfile`
  - https://hub.docker.com/_/golang
  - https://github.com/golangci/golangci-lint/releases
- `terraform/env/terragrunt.hcl`
  - `aws_version`: https://registry.terraform.io/providers/hashicorp/aws/latest/docs
  - `terragrunt_version`: https://formulae.brew.sh/formula/terragrunt

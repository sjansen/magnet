locals {
  aws_version        = "~> 3.27"
  terragrunt_version = "~> 0.28.2"

  env  = path_relative_to_include()
  proj = "magnet"
  region = {
    production = "us-west-2"
    staging    = "us-east-2"
  }[local.env]

  rs_config = find_in_parent_folders("terragrunt-local.json", "")
  rs_prefix = local.rs_config == "" ? local.proj : jsondecode(file(local.rs_config)).rs_prefix
}

generate "locals-provider" {
  path      = "locals-generated.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
locals {
  env  = "${local.env}"
  proj = "${local.proj}"
}
EOF
}

generate "providers" {
  path      = "providers-generated.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "aws" {
  region  = "${local.region}"
}

provider "aws" {
  alias   = "us-east-1"
  region  = "us-east-1"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "${local.aws_version}"
    }
  }
}
EOF
}

remote_state {
  backend = "s3"
  config = {
    region         = local.region
    dynamodb_table = "terraform"
    bucket         = "${local.rs_prefix}-terraform-${local.region}"
    key            = "${local.proj}/${local.env}.tfstate"
    encrypt        = true
  }
}

terragrunt_version_constraint = local.terragrunt_version

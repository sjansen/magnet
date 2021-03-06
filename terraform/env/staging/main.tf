module "app" {
  source = "../../modules/app"
  tags   = local.tags

  cloudwatch-retention = 30
  dns-name             = var.dns-name
  dns-zone             = var.dns-zone
  ssm-prefix           = module.bootstrap.ssm-prefix
  repo-urls            = module.bootstrap.repo-urls

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}

module "bootstrap" {
  source = "../../modules/bootstrap"
  tags   = local.tags

  dns-name = var.dns-name

  providers = {
    aws = aws
  }
}

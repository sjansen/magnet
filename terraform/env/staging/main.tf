module "app" {
  source = "../../modules/app"
  tags   = local.tags

  cloudwatch-retention = 30
  dns-name             = var.dns-name
  dns-zone             = var.dns-zone
  ssm-prefix           = module.bootstrap.ssm-prefix
  webui-repo-url       = module.bootstrap.webui-repo-url

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

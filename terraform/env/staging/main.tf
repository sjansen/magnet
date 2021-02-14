module "app" {
  source = "../../modules/app"

  dns_name = var.dns_name
  dns_zone = var.dns_zone
  tags     = local.tags

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}

locals {
  app_cloudwatch_prefix = "/aws/lambda/${local.app_function_name}"
  app_function_name     = replace(var.dns_name, "/[^-_a-zA-Z0-9]+/", "_")
  app_gateway_name      = var.dns_name
  app_iam_role_name     = var.dns_name
  sessions_table        = "${var.dns_name}-sessions"
  ssm_prefix            = var.dns_name

  cloudwatch_retention = 90

  ecr_repository_name = var.dns_name
}
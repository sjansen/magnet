locals {
  api_gateway_name  = var.dns_name
  web_function_name = replace(var.dns_name, "/[^-_a-zA-Z0-9]+/", "_")
  web_iam_role_name = var.dns_name

  bucket_prefix = replace(var.dns_name, "/[^-_a-zA-Z0-9]+/", "-")

  cloudwatch_apigw_prefix  = "/aws/apigateway/${local.web_function_name}"
  cloudwatch_lambda_prefix = "/aws/lambda/${local.web_function_name}"
  cloudwatch_retention     = 90

  ecr_repository_name = var.dns_name

  sessions_table = "${var.dns_name}-sessions"
  ssm_prefix     = var.dns_name
}

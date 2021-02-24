resource "aws_cloudwatch_log_group" "web" {
  name = local.cloudwatch_lambda_prefix
  tags = var.tags

  retention_in_days = local.cloudwatch_retention
}

resource "aws_cloudwatch_log_group" "apigw" {
  name = local.cloudwatch_apigw_prefix
  tags = var.tags

  retention_in_days = local.cloudwatch_retention
}

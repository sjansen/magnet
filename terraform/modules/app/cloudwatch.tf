resource "aws_cloudwatch_log_group" "apigw" {
  name = "/aws/apigateway/${local.dns-name-underscored}-webui"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}

resource "aws_cloudwatch_log_group" "x" {
  for_each = local.fn-names
  name     = "/aws/lambda/${each.value}"
  tags     = var.tags

  retention_in_days = var.cloudwatch-retention
}

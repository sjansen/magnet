resource "aws_cloudwatch_log_group" "apigw" {
  name = "/aws/apigateway/${local.webui-fn-name}"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}

resource "aws_cloudwatch_log_group" "move" {
  name = "/aws/lambda/${local.move-fn-name}"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}

resource "aws_cloudwatch_log_group" "webui" {
  name = "/aws/lambda/${local.webui-fn-name}"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}

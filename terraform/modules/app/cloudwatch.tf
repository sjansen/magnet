resource "aws_cloudwatch_log_group" "app" {
  name = local.app_cloudwatch_prefix
  tags = var.tags

  retention_in_days = local.cloudwatch_retention
}

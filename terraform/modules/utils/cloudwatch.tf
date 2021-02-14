resource "aws_cloudwatch_log_group" "codebuild" {
  name = local.cloudwatch_prefix
  tags = var.tags

  retention_in_days = local.cloudwatch_retention
}

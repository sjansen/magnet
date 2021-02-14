locals {
  cloudwatch_prefix    = "/aws/codebuild/${local.codebuild_name}"
  cloudwatch_retention = 90

  codebuild_name   = "${local.codebuild_prefix}-promote-docker-image"
  codebuild_prefix = replace(var.dns_name, "/[^-_a-zA-Z0-9]+/", "_")
  codebuild_role   = "${var.dns_name}-codebuild-service-role"
}

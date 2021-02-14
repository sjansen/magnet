locals {
  docker_dst_image    = "${local.docker_dst_registry}/${local.ecr_dst["image"]}:latest"
  docker_dst_registry = "${local.ecr_dst["account"]}.dkr.ecr.${local.ecr_dst["region"]}.amazonaws.com"
  docker_src_image    = "${local.docker_src_registry}/${local.ecr_src["image"]}:latest"
  docker_src_registry = "${local.ecr_src["account"]}.dkr.ecr.${local.ecr_src["region"]}.amazonaws.com"

  ecr_dst = regex("^arn:aws:ecr:(?P<region>[-a-z0-9]+):(?P<account>[0-9]+):repository/(?P<image>.+)$", var.ecr_dst_arn)
  ecr_src = regex("^arn:aws:ecr:(?P<region>[-a-z0-9]+):(?P<account>[0-9]+):repository/(?P<image>.+)$", var.ecr_src_arn)
}

resource "aws_codebuild_project" "promote-docker-image" {
  name           = "${local.codebuild_prefix}-promote-docker-image"
  build_timeout  = "5"
  queued_timeout = "5"
  service_role   = aws_iam_role.codebuild.arn
  tags           = var.tags

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:5.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true

    environment_variable {
      name  = "SOME_KEY1"
      value = "SOME_VALUE1"
    }
  }

  source {
    type      = "NO_SOURCE"
    buildspec = <<EOF
version: 0.2

phases:
  build:
    commands:
      - "aws ecr get-login-password --region ${local.ecr_src["region"]} | docker login --username AWS --password-stdin ${local.docker_src_registry}" 
      - "docker pull ${local.docker_src_image}"
      - "docker tag ${local.docker_src_image} ${local.docker_dst_image}"
      - "aws ecr get-login-password --region ${local.ecr_dst["region"]} | docker login --username AWS --password-stdin ${local.docker_dst_registry}"
      - "aws ecr batch-delete-image --region ${local.ecr_dst["region"]} --repository-name ${local.ecr_dst["image"]} --image-ids imageTag=latest"
      - "docker push ${local.docker_dst_image}"
      - "aws lambda update-function-code --region ${local.ecr_dst["region"]} --function-name ${var.function_name} --image-uri ${local.docker_dst_image}"
EOF
  }
}

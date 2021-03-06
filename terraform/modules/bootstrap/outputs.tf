output "docker-registry" {
  value = split("/", aws_ecr_repository.x["webui"].repository_url)[0]
}

output "ssm-prefix" {
  value = local.ssm-prefix
}

output "repo-arns" {
  value = {
    move  = aws_ecr_repository.x["move"].arn
    webui = aws_ecr_repository.x["webui"].arn
  }
}

output "repo-urls" {
  value = {
    move  = aws_ecr_repository.x["move"].repository_url
    webui = aws_ecr_repository.x["webui"].repository_url
  }
}

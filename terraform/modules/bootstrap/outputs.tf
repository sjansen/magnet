output "docker-registry" {
  value = split("/", aws_ecr_repository.webui.repository_url)[0]
}

output "ssm-prefix" {
  value = local.ssm-prefix
}

output "webui-repo-arn" {
  value = aws_ecr_repository.webui.arn
}

output "webui-repo-url" {
  value = aws_ecr_repository.webui.repository_url
}

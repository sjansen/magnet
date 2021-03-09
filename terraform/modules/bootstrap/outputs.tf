output "docker-registry" {
  value = split("/", aws_ecr_repository.x["webui"].repository_url)[0]
}

output "ssm-prefix" {
  value = local.ssm-prefix
}

output "repo-arns" {
  value = {
    for x in local.lambdas : x => aws_ecr_repository.x[x].arn
  }
}

output "repo-urls" {
  value = {
    for x in local.lambdas : x => aws_ecr_repository.x[x].repository_url
  }
}

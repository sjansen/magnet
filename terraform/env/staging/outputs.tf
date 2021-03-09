output "docker-registry" {
  value = module.bootstrap.docker-registry
}

output "media-bucket" {
  value = module.app.media-bucket
}

output "convert-fn-name" {
  value = module.app.fn-names["convert"]
}

output "convert-repo-arn" {
  value = module.bootstrap.repo-arns["convert"]
}

output "convert-repo-url" {
  value = module.bootstrap.repo-urls["convert"]
}

output "move-fn-name" {
  value = module.app.fn-names["move"]
}

output "move-repo-arn" {
  value = module.bootstrap.repo-arns["move"]
}

output "move-repo-url" {
  value = module.bootstrap.repo-urls["move"]
}

output "ssm-prefix" {
  value = "/${module.bootstrap.ssm-prefix}/"
}

output "webui-fn-name" {
  value = module.app.fn-names["webui"]
}

output "webui-repo-arn" {
  value = module.bootstrap.repo-arns["webui"]
}

output "webui-repo-url" {
  value = module.bootstrap.repo-urls["webui"]
}

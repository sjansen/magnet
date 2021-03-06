output "docker-registry" {
  value = module.bootstrap.docker-registry
}

output "media-bucket" {
  value = module.app.media-bucket
}

output "ssm-prefix" {
  value = "/${module.bootstrap.ssm-prefix}/"
}

output "webui-fn-name" {
  value = module.app.webui-fn-name
}

output "webui-repo-arn" {
  value = module.bootstrap.webui-repo-arn
}

output "webui-repo-url" {
  value = module.bootstrap.webui-repo-url
}

output "bucket" {
  value = module.app.bucket
}

output "ecr_arn" {
  value = module.app.ecr_arn
}

output "function-name" {
  value = module.app.function_name
}

output "registry" {
  value = module.app.registry
}

output "repository-arn" {
  value = module.app.repository-arn
}

output "repository-url" {
  value = module.app.repository-url
}

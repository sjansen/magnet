output "bucket" {
  value = aws_s3_bucket.media.id
}

output "ecr_arn" {
  value = aws_ecr_repository.app.arn
}

output "function_arn" {
  value = aws_lambda_function.app.arn
}

output "function_name" {
  value = aws_lambda_function.app.function_name
}

output "registry" {
  value = split("/", aws_ecr_repository.app.repository_url)[0]
}

output "repository-arn" {
  value = aws_ecr_repository.app.arn
}

output "repository-url" {
  value = aws_ecr_repository.app.repository_url
}

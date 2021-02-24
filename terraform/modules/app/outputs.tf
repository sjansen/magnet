output "bucket" {
  value = aws_s3_bucket.media.id
}

output "ecr_arn" {
  value = aws_ecr_repository.web.arn
}

output "function_arn" {
  value = aws_lambda_function.web.arn
}

output "function_name" {
  value = aws_lambda_function.web.function_name
}

output "registry" {
  value = split("/", aws_ecr_repository.web.repository_url)[0]
}

output "repository-arn" {
  value = aws_ecr_repository.web.arn
}

output "repository-url" {
  value = aws_ecr_repository.web.repository_url
}

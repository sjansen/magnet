output "media-bucket" {
  value = aws_s3_bucket.media.id
}

output "webui-fn-arn" {
  value = aws_lambda_function.webui.arn
}

output "webui-fn-name" {
  value = aws_lambda_function.webui.function_name
}

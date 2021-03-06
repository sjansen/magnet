output "fn-arns" {
  value = {
    move  = aws_lambda_function.move.arn
    webui = aws_lambda_function.webui.arn
  }
}

output "fn-names" {
  value = {
    move  = aws_lambda_function.move.function_name
    webui = aws_lambda_function.webui.function_name
  }
}

output "media-bucket" {
  value = aws_s3_bucket.media.id
}

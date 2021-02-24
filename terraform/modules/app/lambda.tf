resource "aws_lambda_function" "web" {
  image_uri    = "${aws_ecr_repository.web.repository_url}:latest"
  package_type = "Image"
  tags         = var.tags

  function_name = local.web_function_name
  memory_size   = 128
  publish       = true
  role          = aws_iam_role.web.arn
  timeout       = 15

  environment {
    variables = {
      MAGNET_SAML_CERTIFICATE  = aws_ssm_parameter.SAML_METADATA_URL.name
      MAGNET_SAML_METADATA_URL = aws_ssm_parameter.SAML_METADATA_URL.name
      MAGNET_SAML_PRIVATE_KEY  = aws_ssm_parameter.SAML_PRIVATE_KEY.name
      MAGNET_SESSION_TABLE     = aws_dynamodb_table.sessions.name
      MAGNET_SSM_PREFIX        = "/${local.ssm_prefix}/"
      MAGNET_BUCKET            = aws_s3_bucket.media.id,
      MAGNET_URL               = "https://${var.dns_name}/"
    }
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.web,
  ]
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.web.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.web.execution_arn}/*/*"
}

resource "aws_lambda_function" "convert" {
  image_uri    = "${var.repo-urls["convert"]}:latest"
  package_type = "Image"
  tags         = var.tags

  function_name = local.fn-names["convert"]
  memory_size   = 128
  publish       = true
  role          = aws_iam_role.x["convert"].arn
  timeout       = 15

  environment {
    variables = {
      MAGNET_BUCKET = aws_s3_bucket.media.id,
    }
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.x["convert"],
  ]
}

resource "aws_lambda_function" "move" {
  image_uri    = "${var.repo-urls["move"]}:latest"
  package_type = "Image"
  tags         = var.tags

  function_name = local.fn-names["move"]
  memory_size   = 128
  publish       = true
  role          = aws_iam_role.x["move"].arn
  timeout       = 15

  environment {
    variables = {
      MAGNET_BUCKET = aws_s3_bucket.media.id,
    }
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.x["move"],
  ]
}

resource "aws_lambda_function" "webui" {
  image_uri    = "${var.repo-urls["webui"]}:latest"
  package_type = "Image"
  tags         = var.tags

  function_name = local.fn-names["webui"]
  memory_size   = 128
  publish       = true
  role          = aws_iam_role.x["webui"].arn
  timeout       = 15

  environment {
    variables = {
      MAGNET_BUCKET                 = aws_s3_bucket.media.id,
      MAGNET_CLOUDFRONT_KEY_ID      = "ssm"
      MAGNET_CLOUDFRONT_PRIVATE_KEY = "ssm"
      MAGNET_APP_URL                = "https://${var.dns-name}/"
      MAGNET_SAML_CERTIFICATE       = "ssm"
      MAGNET_SAML_METADATA_URL      = "ssm"
      MAGNET_SAML_PRIVATE_KEY       = "ssm"
      MAGNET_SESSION_TABLE          = aws_dynamodb_table.sessions.name
      MAGNET_SSM_PREFIX             = "/${var.ssm-prefix}/"
    }
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.x["webui"],
  ]
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.webui.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.webui.execution_arn}/*/*"
}

resource "aws_lambda_permission" "lambda" {
  statement_id  = "AllowExecutionFromS3"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.move.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.media.arn
}

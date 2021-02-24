resource "aws_api_gateway_account" "apigw" {
  cloudwatch_role_arn = aws_iam_role.apigw.arn
}

resource "aws_api_gateway_base_path_mapping" "default" {
  api_id      = aws_api_gateway_rest_api.web.id
  stage_name  = aws_api_gateway_stage.default.stage_name
  domain_name = aws_api_gateway_domain_name.web.domain_name
}

resource "aws_api_gateway_domain_name" "web" {
  domain_name              = var.dns_name
  regional_certificate_arn = aws_acm_certificate_validation.apigw-cert.certificate_arn
  security_policy          = "TLS_1_2"
  tags                     = var.tags

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "default" {
  rest_api_id = join("", aws_api_gateway_rest_api.web.*.id)
  stage_name  = "default"

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    aws_api_gateway_integration.lambda,
    aws_api_gateway_integration.lambda_root,
  ]
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.web.id
  resource_id = aws_api_gateway_method.proxy.resource_id
  http_method = aws_api_gateway_method.proxy.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.web.invoke_arn
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = aws_api_gateway_rest_api.web.id
  resource_id = aws_api_gateway_method.proxy_root.resource_id
  http_method = aws_api_gateway_method.proxy_root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.web.invoke_arn
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = aws_api_gateway_rest_api.web.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id   = aws_api_gateway_rest_api.web.id
  resource_id   = aws_api_gateway_rest_api.web.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method_settings" "default" {
  rest_api_id = aws_api_gateway_rest_api.web.id
  stage_name  = aws_api_gateway_stage.default.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "ERROR"
  }
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.web.id
  parent_id   = aws_api_gateway_rest_api.web.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_rest_api" "web" {
  name = local.api_gateway_name
  tags = var.tags

  minimum_compression_size = 65536
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_stage" "default" {
  rest_api_id   = aws_api_gateway_rest_api.web.id
  deployment_id = aws_api_gateway_deployment.default.id
  stage_name    = "default"

  xray_tracing_enabled = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.apigw.arn
    format          = <<EOT
$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] "$context.httpMethod $context.resourcePath $context.protocol" $context.status $context.responseLength $context.requestId
EOT
  }

  depends_on = [aws_cloudwatch_log_group.apigw]
}

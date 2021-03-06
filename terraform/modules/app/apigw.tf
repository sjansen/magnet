resource "aws_api_gateway_account" "apigw" {
  cloudwatch_role_arn = aws_iam_role.apigw.arn
}

resource "aws_api_gateway_base_path_mapping" "webui" {
  api_id      = aws_api_gateway_rest_api.webui.id
  stage_name  = aws_api_gateway_stage.webui.stage_name
  domain_name = aws_api_gateway_domain_name.webui.domain_name
}

resource "aws_api_gateway_domain_name" "webui" {
  domain_name              = var.dns-name
  regional_certificate_arn = aws_acm_certificate_validation.apigw-cert.certificate_arn
  security_policy          = "TLS_1_2"
  tags                     = var.tags

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "webui" {
  rest_api_id = join("", aws_api_gateway_rest_api.webui.*.id)

  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    redeployment = sha1(jsonencode([
      aws_lambda_function.webui.version,
      aws_api_gateway_integration.webui.id,
      aws_api_gateway_integration.webui_root.id,
      aws_api_gateway_method.webui.id,
      aws_api_gateway_method.webui_root.id,
      aws_api_gateway_resource.webui.id,
    ]))
  }

  depends_on = [
    aws_api_gateway_integration.webui,
    aws_api_gateway_integration.webui_root,
  ]
}

resource "aws_api_gateway_integration" "webui" {
  rest_api_id = aws_api_gateway_rest_api.webui.id
  resource_id = aws_api_gateway_method.webui.resource_id
  http_method = aws_api_gateway_method.webui.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.webui.invoke_arn
}

resource "aws_api_gateway_integration" "webui_root" {
  rest_api_id = aws_api_gateway_rest_api.webui.id
  resource_id = aws_api_gateway_method.webui_root.resource_id
  http_method = aws_api_gateway_method.webui_root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.webui.invoke_arn
}

resource "aws_api_gateway_method" "webui" {
  rest_api_id   = aws_api_gateway_rest_api.webui.id
  resource_id   = aws_api_gateway_resource.webui.id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method" "webui_root" {
  rest_api_id   = aws_api_gateway_rest_api.webui.id
  resource_id   = aws_api_gateway_rest_api.webui.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method_settings" "webui" {
  rest_api_id = aws_api_gateway_rest_api.webui.id
  stage_name  = aws_api_gateway_stage.webui.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "ERROR"
  }
}

resource "aws_api_gateway_resource" "webui" {
  rest_api_id = aws_api_gateway_rest_api.webui.id
  parent_id   = aws_api_gateway_rest_api.webui.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_rest_api" "webui" {
  name = local.api_gateway_name
  tags = var.tags

  minimum_compression_size = 65536
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_stage" "webui" {
  rest_api_id   = aws_api_gateway_rest_api.webui.id
  deployment_id = aws_api_gateway_deployment.webui.id
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

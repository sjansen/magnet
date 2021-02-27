data "aws_iam_policy_document" "AssumeRole-apigw" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "AssumeRole-lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "web" {
  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:Query",
      "dynamodb:UpdateItem",
    ]
    resources = [aws_dynamodb_table.sessions.arn]
  }
  statement {
    actions = [
      "s3:GetObject",
      "s3:ListBucket",
    ]
    resources = [
      aws_s3_bucket.media.arn,
      "${aws_s3_bucket.media.arn}/*",
    ]
  }
  statement {
    actions = [
      "s3:PutObject",
      "s3:PutObjectAcl",
    ]
    resources = ["${aws_s3_bucket.media.arn}/inbox/*"]
  }
  statement {
    actions   = ["ssm:GetParameters"]
    resources = ["arn:aws:ssm:*:*:parameter/${local.ssm_prefix}/*"]
  }
}

resource "aws_iam_policy" "web" {
  name   = "all-the-things"
  path   = "/"
  policy = data.aws_iam_policy_document.web.json
}

resource "aws_iam_role" "apigw" {
  name = "APIGateway"
  tags = var.tags

  assume_role_policy = data.aws_iam_policy_document.AssumeRole-apigw.json
}

resource "aws_iam_role" "web" {
  name = local.web_iam_role_name
  tags = var.tags

  assume_role_policy = data.aws_iam_policy_document.AssumeRole-lambda.json
}

resource "aws_iam_role_policy_attachment" "apigw" {
  role       = aws_iam_role.apigw.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

resource "aws_iam_role_policy_attachment" "web" {
  role       = aws_iam_role.web.name
  policy_arn = aws_iam_policy.web.arn
}

resource "aws_iam_role_policy_attachment" "web-logs" {
  role       = aws_iam_role.web.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "web-xray" {
  role       = aws_iam_role.web.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

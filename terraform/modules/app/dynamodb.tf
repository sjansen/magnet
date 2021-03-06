resource "aws_dynamodb_table" "sessions" {
  name = "${var.dns-name}-sessions"
  tags = var.tags

  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "token"

  attribute {
    name = "token"
    type = "S"
  }

  server_side_encryption {
    enabled = true
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }
}

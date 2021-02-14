resource "aws_s3_bucket" "logs" {
  bucket = "${var.dns_name}-logs"
  acl    = "log-delivery-write"

  tags = var.tags

  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 90
    }
    noncurrent_version_expiration {
      days = 30
    }
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "media" {
  bucket = "${local.bucket_prefix}-media"
  tags   = var.tags

  acceleration_status = "Enabled"
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["POST"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 300
  }
  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    noncurrent_version_expiration {
      days = 30
    }
  }
  lifecycle_rule {
    id                                     = "clean inbox"
    enabled                                = true
    prefix                                 = "inbox/"
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 7
    }
    noncurrent_version_expiration {
      days = 7
    }
  }
  lifecycle_rule {
    id                                     = "clean review"
    enabled                                = true
    prefix                                 = "review/"
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 45
    }
    noncurrent_version_expiration {
      days = 45
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

resource "aws_s3_bucket" "logs" {
  bucket = "${local.bucket_prefix}-logs"
  acl    = "log-delivery-write"
  tags   = var.tags

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

resource "aws_s3_bucket_notification" "media" {
  bucket = aws_s3_bucket.media.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.move.arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "inbox/"
  }

  depends_on = [aws_lambda_permission.lambda]
}

resource "aws_s3_bucket_object" "favicon" {
  bucket       = aws_s3_bucket.media.id
  key          = "magnet/icons/favicon.ico"
  content_type = "image/x-icon"
  etag         = filemd5("${path.module}/icons/favicon.ico")
  source       = "${path.module}/icons/favicon.ico"
}

resource "aws_s3_bucket_object" "icons" {
  for_each = fileset(path.module, "icons/*.svg")

  bucket       = aws_s3_bucket.media.id
  key          = "magnet/${each.key}"
  content_type = "image/svg+xml"
  etag         = filemd5("${path.module}/${each.key}")
  source       = "${path.module}/${each.key}"
}

resource "aws_s3_bucket_object" "prefixes" {
  for_each = toset(["inbox/", "review/", "static/"])

  bucket = aws_s3_bucket.media.id
  key    = each.key
}

resource "aws_s3_bucket_policy" "media" {
  bucket = aws_s3_bucket.media.id
  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [{
    "Effect":"Allow",
    "Action": "s3:GetObject",
    "Principal": {
      "AWS": "${aws_cloudfront_origin_access_identity.cdn.iam_arn}"
    },
    "Resource": [
      "${aws_s3_bucket.media.arn}/magnet/icons/*",
      "${aws_s3_bucket.media.arn}/review/*",
      "${aws_s3_bucket.media.arn}/static/*"
    ]
  }]
}
EOF
}

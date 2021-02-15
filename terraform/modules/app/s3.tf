resource "aws_s3_bucket" "logs" {
  bucket = "${var.dns_name}-logs"
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

resource "aws_s3_bucket" "media" {
  bucket = "${var.dns_name}-media"
  tags   = var.tags

  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
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

resource "aws_s3_bucket_object" "icons" {
  for_each = fileset(path.module, "icons/*.svg")

  bucket = aws_s3_bucket.media.id
  key    = each.key

  content_type = "image/svg+xml"
  etag         = filemd5("${path.module}/${each.key}")
  source       = "${path.module}/${each.key}"
}

resource "aws_s3_bucket_object" "prefixes" {
  for_each = toset(["inbox/", "media/", "review/"])

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
    "Resource": "${aws_s3_bucket.media.arn}/icons/*"
  },{
    "Effect":"Allow",
    "Action": "s3:GetObject",
    "Principal": {
      "AWS": "${aws_cloudfront_origin_access_identity.cdn.iam_arn}"
    },
    "Resource": "${aws_s3_bucket.media.arn}/media/*"
  }]
}
EOF
}

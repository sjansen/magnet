resource "aws_ssm_parameter" "CLOUDFRONT_KEY_ID" {
  name        = "/${local.ssm-prefix}/CLOUDFRONT_KEY_ID"
  description = "$MAGNET_CLOUDFRONT_KEY_ID"
  type        = "String"
  value       = "invalid"
  overwrite   = false
  tags        = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "CLOUDFRONT_PRIVATE_KEY" {
  name        = "/${local.ssm-prefix}/CLOUDFRONT_PRIVATE_KEY"
  description = "$MAGNET_CLOUDFRONT_PRIVATE_KEY"
  type        = "SecureString"
  value       = "invalid"
  overwrite   = false
  tags        = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "SAML_CERTIFICATE" {
  name        = "/${local.ssm-prefix}/SAML_CERTIFICATE"
  description = "$MAGNET_SAML_CERTIFICATE"
  type        = "String"
  value       = "invalid"
  overwrite   = false
  tags        = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "SAML_METADATA_URL" {
  name        = "/${local.ssm-prefix}/SAML_METADATA_URL"
  description = "$MAGNET_SAML_METADATA_URL"
  type        = "String"
  value       = "invalid"
  overwrite   = false
  tags        = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "SAML_PRIVATE_KEY" {
  name        = "/${local.ssm-prefix}/SAML_PRIVATE_KEY"
  description = "$MAGNET_SAML_PRIVATE_KEY"
  type        = "SecureString"
  value       = "invalid"
  overwrite   = false
  tags        = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

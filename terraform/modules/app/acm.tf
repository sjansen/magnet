resource "aws_acm_certificate" "apigw-cert" {
  domain_name       = var.dns-name
  validation_method = "DNS"
  tags              = var.tags
}


resource "aws_acm_certificate" "cloudfront-cert" {
  provider          = aws.us-east-1
  domain_name       = var.dns-name
  validation_method = "DNS"
  tags              = var.tags
}


resource "aws_acm_certificate_validation" "apigw-cert" {
  certificate_arn         = aws_acm_certificate.apigw-cert.arn
  validation_record_fqdns = [for record in aws_route53_record.apigw-cert : record.fqdn]
}


resource "aws_acm_certificate_validation" "cloudfront-cert" {
  provider                = aws.us-east-1
  certificate_arn         = aws_acm_certificate.cloudfront-cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cloudfront-cert : record.fqdn]
}

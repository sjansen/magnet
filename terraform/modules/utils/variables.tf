variable "dns_name" {
  type = string
}

variable "ecr_dst_arn" {
  type        = string
  description = "e.g. arn:aws:ecr:us-east-1:123456789012:repository/magnet-prod"
}

variable "ecr_src_arn" {
  type        = string
  description = "e.g. arn:aws:ecr:us-east-1:123456789012:repository/magnet-devel"
}

variable "function_arn" {
  type = string
}

variable "function_name" {
  type = string
}

variable "tags" {
  type = map(string)
}

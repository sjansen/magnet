variable "cloudwatch-retention" {
  type = number
}

variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "repo-urls" {
  type = map(string)
}

variable "ssm-prefix" {
  type = string
}

variable "tags" {
  type = map(string)
}

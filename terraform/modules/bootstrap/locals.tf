locals {
  bucket_prefix = local.dns-name-dashed

  dns-name-dashed      = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "-")
  dns-name-underscored = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "_")

  lambdas = toset(["convert", "move", "webui"])

  ssm-prefix = var.dns-name
}

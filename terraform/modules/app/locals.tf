locals {
  api_gateway_name = var.dns-name

  bucket_prefix = local.dns-name-dashed

  dns-name-dashed      = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "-")
  dns-name-underscored = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "_")

  fn-names = {
    for x in local.lambdas : x => "${local.dns-name-underscored}-${x}"
  }

  lambdas = toset(["convert", "move", "webui"])
}

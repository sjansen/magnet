locals {
  api_gateway_name = var.dns-name

  bucket_prefix = local.dns-name-dashed

  dns-name-dashed      = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "-")
  dns-name-underscored = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "_")

  move-fn-name  = "${local.dns-name-underscored}-move"
  webui-fn-name = "${local.dns-name-underscored}-webui"
}

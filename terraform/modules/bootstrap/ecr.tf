resource "aws_ecr_repository" "webui" {
  name = "${var.dns-name}-webui"
  tags = var.tags

  image_tag_mutability = "IMMUTABLE"
}

resource "aws_ecr_repository_policy" "webui" {
  repository = aws_ecr_repository.webui.name
  policy     = data.aws_iam_policy_document.webui-ecr.json
}

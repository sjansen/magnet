data "aws_iam_policy_document" "ecr-web" {
  statement {
    sid = "LambdaECRImageRetrievalPolicy"
    actions = [
      "ecr:BatchGetImage",
      "ecr:DeleteRepositoryPolicy",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetRepositoryPolicy",
      "ecr:SetRepositoryPolicy",
    ]
    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

resource "aws_ecr_repository" "web" {
  name = local.ecr_repository_name
  tags = var.tags

  image_tag_mutability = "IMMUTABLE"
}

resource "aws_ecr_repository_policy" "web" {
  repository = aws_ecr_repository.web.name
  policy     = data.aws_iam_policy_document.ecr-web.json
}

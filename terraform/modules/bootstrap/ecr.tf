resource "aws_ecr_lifecycle_policy" "webui" {
  repository = aws_ecr_repository.webui.name
  policy     = <<EOF
{
    "rules": [{
        "rulePriority": 10,
        "description": "Expire untagged images after 3 days",
        "selection": {
            "tagStatus": "untagged",
            "countType": "sinceImagePushed",
            "countUnit": "days",
            "countNumber": 3
        },
        "action": {
            "type": "expire"
        }
    }, {
        "rulePriority": 100,
        "description": "Keep last 3 tagged images",
        "selection": {
            "tagStatus": "any",
            "countType": "imageCountMoreThan",
            "countNumber": 3
        },
        "action": {
            "type": "expire"
        }
    }]
}
EOF
}

resource "aws_ecr_repository" "webui" {
  name = "${var.dns-name}-webui"
  tags = var.tags

  image_tag_mutability = "IMMUTABLE"
}

resource "aws_ecr_repository_policy" "webui" {
  repository = aws_ecr_repository.webui.name
  policy     = data.aws_iam_policy_document.webui-ecr.json
}

resource "aws_ecr_lifecycle_policy" "x" {
  for_each = local.lambdas

  repository = aws_ecr_repository.x[each.key].name
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

resource "aws_ecr_repository" "x" {
  for_each = local.lambdas

  name = "${var.dns-name}-${each.key}"
  tags = var.tags

  image_tag_mutability = "IMMUTABLE"
}

resource "aws_ecr_repository_policy" "x" {
  for_each = local.lambdas

  repository = aws_ecr_repository.x[each.key].name
  policy     = data.aws_iam_policy_document.ecr-lambda.json
}

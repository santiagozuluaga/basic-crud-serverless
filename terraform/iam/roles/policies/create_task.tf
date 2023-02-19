resource "aws_iam_policy" "create_task" {

  name   = "create-task"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "dynamodb:PutItem"
      ],
      "Resource": [
          "arn:aws:dynamodb:${var.inputs.config.region}:${var.inputs.config.account_id}:table/task"
      ]
    }
  ]
}
EOF
}

output "create_task" {
  value = aws_iam_policy.create_task
}

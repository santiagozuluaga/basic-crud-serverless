resource "aws_iam_role" "iam_lambda_role" {
  name = var.name

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_iam_lambda_policy_to_iam_lambda_role" {
  count = length(var.policy_arns)

  role       = aws_iam_role.iam_lambda_role.name
  policy_arn = var.policy_arns[count.index]
}

output "role" {
  value = aws_iam_role.iam_lambda_role
}

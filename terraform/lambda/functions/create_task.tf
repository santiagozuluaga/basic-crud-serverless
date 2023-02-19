resource "aws_lambda_function" "create_task" {
  function_name    = "create_task"
  handler          = "create-task"
  runtime          = "go1.x"
  role             = var.inputs.modules.iam.create_task.arn
  filename         = "../bin/create-task.zip"
  source_code_hash = sha256(filebase64("../bin/create-task.zip"))
  memory_size      = 128
  timeout          = 5
}

output "create_task" {
  value = aws_lambda_function.create_task
}

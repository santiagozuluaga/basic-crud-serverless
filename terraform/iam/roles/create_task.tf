module "create_task" {
  source = "../../shared/role"

  name = "create-task"

  policy_arns = [
    module.policies.create_task.arn,
    module.policies.cloudwatch_logs.arn
  ]
}

output "create_task" {
  value = module.create_task.role
}
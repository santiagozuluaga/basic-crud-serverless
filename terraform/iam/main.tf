module "roles" {
  source = "./roles"

  inputs = var.inputs
}

output "roles" {
  value = module.roles
}
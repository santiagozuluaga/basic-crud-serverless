module "functions" {
  source = "./functions"

  inputs = var.inputs
}

output "functions" {
  value = module.functions
}
module "tables" {
  source = "./tables"

  inputs = var.inputs
}

output "tables" {
  value = module.tables
}
resource "terraform_data" "env_guard" {
  lifecycle {
    precondition {
      condition     = var.env == terraform.workspace
      error_message = "El workspace '${terraform.workspace}' no coincide con el env '${var.env}'"
    }
  }
}

module "network" {
  source = "./modules/network"

  name = "tero-${var.env}"
  vpc_cidr = var.vpc_cidr
  public_subnet_cidrs = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}
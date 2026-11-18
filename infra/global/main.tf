locals {
  services = ["catalog", "cart", "orders", "checkout", "ui", "admin"]
}

module "ecr" {
  source = "../modules/ecr"

  repository_names = toset([for s in local.services : "tero-${s}"])
  force_delete     = true
}
output "vpc_id" {
  description = "ID de la VPC del ambiente"
  value       = module.network.vpc_id
}

output "public_subnet_ids" {
  description = "IDs de las subredes publicas"
  value       = module.network.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs de las subredes privadas"
  value       = module.network.private_subnet_ids
}
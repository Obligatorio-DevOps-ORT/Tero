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

output "nat_public_ip" {
  description = "IP publica de salida de las subredes privadas (null si el NAT esta apagado)"
  value       = module.network.nat_public_ip
}

output "ecs_cluster_name" {
  description = "Nombre del cluster ECS del ambiente"
  value       = module.ecs_cluster.cluster_name
}

output "ecs_cluster_arn" {
  description = "ARN del cluster ECS del ambiente"
  value       = module.ecs_cluster.cluster_arn
}

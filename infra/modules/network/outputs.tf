output "vpc_id" {
  description = "ID de la VPC"
  value       = aws_vpc.this.id
}

output "vpc_cidr" {
  description = "Rango de IPs de la VPC, para las reglas del security groups"
  value       = aws_vpc.this.cidr_block
}

output "public_subnet_ids" {
  description = "IDs de las subredes publicas, para el ALB"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs de las subredes privadas, para el ECS, RDS y Redis"
  value       = aws_subnet.private[*].id
}

output "alb_sg_id" {
  description = "ID del security group del ALB"
  value       = aws_security_group.alb.id
}

output "ecs_sg_id" {
  description = "ID del security group de los servicios ECS"
  value       = aws_security_group.ecs.id
}

output "rds_sg_id" {
  description = "ID del security group de Postgres (RDS)"
  value       = aws_security_group.rds.id
}

output "redis_sg_id" {
  description = "ID del security group de Redis (ElastiCache)"
  value       = aws_security_group.redis.id
}

output "nat_public_ip" {
  description = "IP publica del NAT Gateway (null si enable_nat_gateway es false)"
  value       = one(aws_eip.nat[*].public_ip)
}

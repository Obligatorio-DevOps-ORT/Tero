output "cluster_id" {
  description = "ID del cluster ECS"
  value       = aws_ecs_cluster.this.id
}

output "cluster_arn" {
  description = "ARN del cluster ECS, lo usan los services y las politicas IAM"
  value       = aws_ecs_cluster.this.arn
}

output "cluster_name" {
  description = "Nombre del cluster ECS, lo usa el autoscaling, las alarmas de CloudWatch y el deploy del pipeline"
  value       = aws_ecs_cluster.this.name
}
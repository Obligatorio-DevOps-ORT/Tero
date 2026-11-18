output "repository_urls" {
  description = "URL de cada repositorio por nombre, para el docker push del CI y las task definitions de ECS"
  value       = { for name, repo in aws_ecr_repository.this : name => repo.repository_url }
}
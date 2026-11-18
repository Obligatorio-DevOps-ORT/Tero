output "ecr_repository_urls" {
  description = "URL de cada repositorio ECR por nombre, compartidos por todos los ambientes"
  value       = module.ecr.repository_urls
}
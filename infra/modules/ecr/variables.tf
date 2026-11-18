variable "repository_names" {
  description = "Nombres de los repositorios ECR a crear, uno por servicio"
  type        = set(string)
}

variable "image_retention_count" {
  description = "Cantidad de imagenes que se guardan por repositorio"
  type        = number
  default     = 5

  validation {
    condition     = var.image_retention_count >= 1
    error_message = "image_retention_count tiene que ser al menos 1"
  }
}

variable "force_delete" {
  description = "Permite borrar repositorios que tenga imagenes (solo para entornos efimeros)"
  type        = bool
  default     = false
}
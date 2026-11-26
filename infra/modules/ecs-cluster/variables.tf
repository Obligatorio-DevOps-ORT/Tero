variable "name" {
  description = "Nombre del cluster"
  type        = string
}

variable "enable_container_insights" {
  description = "Activar Container Insights (metrica por task y red). Activarlo solo en prod"
  type        = bool
  default     = false
}

variable "fargate_base" {
  description = "Cantidad minima de tas en Fargate on-demand. El resto va a fargate_spot, 0 = todo spot"
  type        = number
  default     = 0
}

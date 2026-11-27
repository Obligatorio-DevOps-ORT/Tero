variable "env" {
  description = "Ambiente (dev, test, o prod), debe coincidir con el workspace"
  type        = string
}

variable "vpc_cidr" {
  description = "Rango de IPs de la VPC del ambiente"
  type        = string
}

variable "public_subnet_cidrs" {
  description = "CIDRs de las subredes publicas"
  type        = list(string)
}

variable "private_subnet_cidrs" {
  description = "CIDRs de las subredes privadas"
  type        = list(string)
}

variable "aws_region" {
  description = "Region AWS"
  type        = string
  default     = "us-east-2"
}

variable "enable_nat_gateway" {
  description = "Crear el NAT Gateway del ambiente (cobra por hora)"
  type        = bool
  default     = false
}

variable "enable_container_insights" {
  description = "Activar Container Insights en el cluster ECS"
  type        = bool
  default     = false
}

variable "fargate_base" {
  description = "Task minimas en FARGATE on-demand, el resto va a FARGATE_SPOT"
  type        = number
  default     = 0
}

variable "name" {
  description = "Prefijo para nombrar los recursos, ej: tero-dev"
  type        = string
}

variable "vpc_cidr" {
  description = "Rango de IPs de la VPC"
  type        = string
}

variable "public_subnet_cidrs" {
  description = "CIDRs de las subredes publicas, una por zona de disponibilidad"
  type        = list(string)
}

variable "private_subnet_cidrs" {
  description = "CIDRs de las subredes privadas, una por zona de disponibilidad"
  type        = list(string)
}

variable "enable_nat_gateway" {
  description = "Crear un NAT Gateway para dar salida a internet a las subredes privadas. Cobra por hora: activarlo solo cuando se use"
  type        = bool
  default     = false
}

variable "app_port" {
  description = "Puerto en el que escucha los contenedores de los servicios"
  type        = number
  default     = 8080
}

variable "db_port" {
  description = "Puerto de Postgres"
  type        = number
  default     = 5432
}

variable "redis_port" {
  description = "Puerto de Redis"
  type        = number
  default     = 6379
}

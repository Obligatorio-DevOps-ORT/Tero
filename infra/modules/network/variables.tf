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
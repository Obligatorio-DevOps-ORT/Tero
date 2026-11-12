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
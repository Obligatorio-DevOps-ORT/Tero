variable "env" {
  type = string
}

resource "terraform_data" "env_guard" {
  lifecycle {
    precondition {
      condition = var.env == terraform.workspace
      error_message = "El workspace '${terraform.workspace}' no coincide con el entorno de enviroment '${var.env}' del tfvars "
    }
  }
}
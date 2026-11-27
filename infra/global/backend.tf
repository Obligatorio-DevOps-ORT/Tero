terraform {
  required_version = ">=1.17"

  backend "s3" {
    bucket       = "tero-tfstate-550842784335"
    key          = "global/terraform.tfstate"
    region       = "us-east-2"
    encrypt      = true
    use_lockfile = true
  }
}

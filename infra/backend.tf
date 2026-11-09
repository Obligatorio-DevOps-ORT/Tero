terraform{
  required_version = ">=1.17"

  backend "s3" {
    bucket = "obligatorio-devops-tfstate-26"
    key = "terraform.tfstate"
    region = "us-east-1"
    encrypt = true
    use_lockfile = true
  }
}
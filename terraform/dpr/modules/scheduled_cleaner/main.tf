terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.27.0"
    }
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
    null_resource = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

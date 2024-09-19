####################################################
# Base
####################################################

terraform {
  required_version = ">= 1.5.7" # 1.5 is required for gitlab

  required_providers {
    aws = {
      source                = "hashicorp/aws"
      version               = "~> 5.0"
    }
  }

  backend "s3" {
    region = "ca-central-1"
    # For production, this value must be overridden via the
    # -backend-config="bucket= ... " command with the production bucket
    bucket         = "heavenly-dragons-terraform-dev-fib720"
    dynamodb_table = "terraform-state-lock"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Managed_By    = "Terraform"
      Env           = var.env
      Project_Topic = "Heavenly Dragons Game Server"
    }
  }
}

locals {
  prefix = "heavenly-dragons"
  suffix = var.env
#   random = "a0cd3657"
}
terraform {
  required_providers {
    aws = {
      version = ">= 2.7.0"
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region                  = var.aws_region
  profile                 = var.aws_profile
}

locals {
  stage = "dev"
}

module "pooling" {
  source = "../modules/pooling"
  stage = local.stage
}
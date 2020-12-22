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

data "aws_caller_identity" "current" {}

locals {
  stage = "dev"
}

module "iam" {
  source = "../../services/iam"
  stage = local.stage
  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region = var.aws_region
}

module "pooling" {
  source = "../../services/pooling"
  stage = local.stage
  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region = var.aws_region
}
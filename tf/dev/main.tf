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

module "conductor" {
  source = "../modules/conductor"
  stage = local.stage
  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region = var.aws_region
}

module "pooling" {
  source = "../modules/pooling"
  stage = local.stage
  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region = var.aws_region
  conductor_queue_arn = module.conductor.this_queue_arn
}
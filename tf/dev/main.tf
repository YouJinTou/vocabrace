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

resource "null_resource" "remove_builds" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "rm -rf builds"
  }
}

module "pooling" {
  source = "../modules/pooling"
  stage = local.stage
  aws_account_id = data.aws_caller_identity.current.account_id
  aws_region = var.aws_region
  depends_on = [null_resource.remove_builds]
}
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
  depends_on = [null_resource.remove_builds]
}
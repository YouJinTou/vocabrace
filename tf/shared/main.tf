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

resource "aws_dynamodb_table" "wordlines_bulgarian" {
  name           = "wordlines_bulgarian"
  read_capacity  = 5
  write_capacity = 1
  hash_key       = "Word"
  attribute {
    name = "Word"
    type = "S"
  }
}

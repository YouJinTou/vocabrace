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

resource "aws_dynamodb_table" "scrabble_bulgarian" {
  name           = "scrabble_bulgarian"
  read_capacity  = 1
  write_capacity = 22
  hash_key       = "Word"
  attribute {
    name = "Word"
    type = "S"
  }
}

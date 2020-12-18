module "pool" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/pool.zip"
  function_name = "${var.stage}_pool"
  handler = "pool"
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
  }
  enable_sns = true
  sns_arn = aws_sns_topic.pools.arn
}

resource "aws_dynamodb_table" "pools" {
  name           = "${var.stage}_pools"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ID"
  attribute {
    name = "ID"
    type = "S"
  }
  ttl {
    attribute_name = "LiveUntil"
    enabled        = true
  }
}

resource "aws_sns_topic" "pools" {
  name           = "${var.stage}_pools"
}
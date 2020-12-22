module "reconnect" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/reconnect.zip"
  function_name = "${var.stage}_reconnect"
  handler = "reconnect"
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
  }
  enable_sns = true
  sns_arn = aws_sns_topic.reconnect.arn
}

resource "aws_sns_topic" "reconnect" {
  name           = "${var.stage}_reconnect"
}

resource "aws_dynamodb_table" "user_pools" {
  name           = "${var.stage}_user_pools"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "UserID"
  attribute {
    name = "UserID"
    type = "S"
  }
  ttl {
    attribute_name = "LiveUntil"
    enabled        = true
  }
}
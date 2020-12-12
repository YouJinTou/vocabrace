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
  enable_streaming = true
  stream_arn = aws_dynamodb_table.tallies.stream_arn
  reserved_concurrent_executions = 1
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

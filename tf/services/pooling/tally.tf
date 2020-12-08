module "tally" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/tally.zip"
  function_name = "${var.stage}_tally"
  handler = "tally"
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
  }
  enable_streaming = true
  stream_arn = aws_dynamodb_table.connections.stream_arn
  reserved_concurrent_executions = 1
}

resource "aws_dynamodb_table" "tallies" {
  name           = "${var.stage}_tallies"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ID"
  stream_enabled = true
  stream_view_type = "NEW_AND_OLD_IMAGES"
  attribute {
    name = "ID"
    type = "S"
  }
}
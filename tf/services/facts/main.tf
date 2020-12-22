resource "aws_dynamodb_table" "facts" {
  name           = "${var.stage}_facts"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ID"
  range_key = "Timestamp"
  attribute {
    name = "ID"
    type = "S"
  }
  attribute {
    name = "Timestamp"
    type = "N"
  }
  stream_enabled = true
  stream_view_type = "NEW_IMAGE"
}

module "broadcast" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/broadcast.zip"
  function_name = "${var.stage}_broadcast"
  handler = "broadcast"
  timeout = 2
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
    ACCOUNT_ID: var.aws_account_id
    IS_SERVERLESS: true
  }
  source_maximum_retries = 5
  enable_streaming = true
  stream_arn = aws_dynamodb_table.facts.stream_arn
}

resource "aws_sns_topic" "facts" {
  name           = "${var.stage}_facts"
}
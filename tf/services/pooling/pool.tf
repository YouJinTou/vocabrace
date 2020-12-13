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
  sqs_sources = [
    {arn: aws_sqs_queue.pools.arn, batch_size: 1}
  ]
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

resource "aws_sqs_queue" "pools" {
  name           = "${var.stage}_pools"
}
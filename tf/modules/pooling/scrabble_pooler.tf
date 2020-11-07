module "scrabble_pooler" {
    source = "../lambda"
    aws_account_id = var.aws_account_id
    aws_region = var.aws_region
    filename = "../payloads/pooler.zip"
    function_name = "${var.stage}_pooler"
    handler = "pooler"
    environment_variables = {
        STAGE: var.stage
        REGION: var.aws_region
        ACCOUNT_ID: var.aws_account_id
        POOL_LIMIT: 3
    }
    function_can_invoke_api_gateway = true
}

resource "aws_sqs_queue" "scrabble_pooler" {
  name                      = "${var.stage}_scrabble_pooler"
  message_retention_seconds = 3600
  visibility_timeout_seconds = 3
}
locals {
  scrabble_language_limit = [
    for pair in setproduct(["bulgarian", "english"], [2, 3, 4]) : {
      language = pair[0]
      players = pair[1]
    }
  ]
  players_by_arn = {
    for q in aws_sqs_queue.scrabble_pooler:
      q.arn => (
        length(regexall("_2_", q.arn)) > 0 ? 2 :
        length(regexall("_3_", q.arn)) > 0 ? 3 :
        4)
  }
}

module "scrabble_pooler" {
    for_each = local.players_by_arn
    source = "../../modules/lambda"
    aws_account_id = var.aws_account_id
    aws_region = var.aws_region
    filename = "../../payloads/pooler.zip"
    function_name = split(":", each.key)[5]
    handler = "pooler"
    environment_variables = {
        STAGE: var.stage
        REGION: var.aws_region
        ACCOUNT_ID: var.aws_account_id
    }
    function_can_invoke_api_gateway = true
    sqs_sources = [
      { arn: each.key, batch_size: each.value }
    ]
    timeout = 5
    reserved_concurrent_executions = 1
    depends_on = [aws_sqs_queue.scrabble_pooler]
}

resource "aws_sqs_queue" "scrabble_pooler" {
  for_each = {for x in local.scrabble_language_limit : "${x.language}_${x.players}" => x}
  name                      = "${var.stage}_scrabble_${each.value.language}_${each.value.players}_pooler"
  message_retention_seconds = 3600
  visibility_timeout_seconds = 5
}
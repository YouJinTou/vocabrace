module "conductor" {
    source = "../lambda"
    aws_account_id = var.aws_account_id
    aws_region = var.aws_region
    filename = "../payloads/conductor.zip"
    function_name = "${var.stage}_conductor"
    handler = "conductor"
    environment_variables = {
        STAGE: var.stage
    }
    function_can_invoke_api_gateway = true
    sqs_can_invoke_function = true
    sqs_source_arn = aws_sqs_queue.conductor.arn
}

resource "aws_sqs_queue" "conductor" {
  name                      = "${var.stage}_conductor"
  message_retention_seconds = 300
}
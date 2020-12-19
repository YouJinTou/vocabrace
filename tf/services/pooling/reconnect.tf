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
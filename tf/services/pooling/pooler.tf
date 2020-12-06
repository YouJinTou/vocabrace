module "pooler" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/pooler.zip"
  function_name = "${var.stage}_pooler"
  handler = "pooler"
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
    ACCOUNT_ID: var.aws_account_id
  }
  function_can_invoke_api_gateway = true
  api_gateway_can_invoke_function = true
  reserved_concurrent_executions = 1
}

resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "$connect"
  target = "integrations/${aws_apigatewayv2_integration.connect.id}"
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "$disconnect"
  target = "integrations/${aws_apigatewayv2_integration.disconnect.id}"
}

resource "aws_apigatewayv2_integration" "connect" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.pooler.this_lambda_function_invoke_arn
}

resource "aws_apigatewayv2_integration" "disconnect" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.pooler.this_lambda_function_invoke_arn
}

resource "aws_dynamodb_table" "waitlist" {
  name           = "${var.stage}_waitlist"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ID"
  attribute {
    name = "ID"
    type = "S"
  }
}
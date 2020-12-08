module "disconnect" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/disconnect.zip"
  function_name = "${var.stage}_disconnect"
  handler = "disconnect"
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
  }
  api_gateway_can_invoke_function = true
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "$disconnect"
  target = "integrations/${aws_apigatewayv2_integration.disconnect.id}"
}

resource "aws_apigatewayv2_integration" "disconnect" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.disconnect.this_lambda_function_invoke_arn
}
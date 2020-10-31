module "connect" {
  source = "../lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "./connect.zip"
  function_name = "${var.stage}_connect"
  handler = "connect"
  environment_variables = {
    STAGE: var.stage
  }
  function_can_invoke_api_gateway = true
  api_gateway_can_invoke_function = true
  api_gateway_source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/$connect"
}

resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "$connect"
  target = "integrations/${aws_apigatewayv2_integration.connect.id}"
}

resource "aws_apigatewayv2_integration" "connect" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.connect.this_lambda_function_invoke_arn
}
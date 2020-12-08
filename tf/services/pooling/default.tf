resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "$default"
  target = "integrations/${aws_apigatewayv2_integration.connect.id}"
}

resource "aws_apigatewayv2_integration" "default" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.connect.this_lambda_function_invoke_arn
}
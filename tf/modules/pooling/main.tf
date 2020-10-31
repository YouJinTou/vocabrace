resource "aws_apigatewayv2_api" "pooling" {
  name                       = "${var.stage}_pooling"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "pooling" {
  api_id = aws_apigatewayv2_api.pooling.id
  name   = "v1"
}

resource "aws_apigatewayv2_deployment" "pooling" {
  api_id      = aws_apigatewayv2_api.pooling.id
  description = "Terraform deployment"

  triggers = {
    redeployment = sha1(join(",", list(
      jsonencode(aws_apigatewayv2_integration.connect),
      jsonencode(aws_apigatewayv2_route.connect),
      jsonencode(aws_apigatewayv2_integration.disconnect),
      jsonencode(aws_apigatewayv2_route.disconnect),
      jsonencode(aws_apigatewayv2_integration.publish),
      jsonencode(aws_apigatewayv2_route.publish),
    )))
  }

  lifecycle {
    create_before_destroy = true
  }
}
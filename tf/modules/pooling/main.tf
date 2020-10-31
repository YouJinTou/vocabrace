resource "aws_apigatewayv2_api" "pooling" {
  name                       = "${var.stage}_pooling"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}
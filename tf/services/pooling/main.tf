resource "aws_apigatewayv2_api" "pooling" {
  name                       = "${var.stage}_pooling"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "pooling" {
  api_id = aws_apigatewayv2_api.pooling.id
  name   = var.stage
}

resource "aws_apigatewayv2_deployment" "pooling" {
  api_id      = aws_apigatewayv2_api.pooling.id
  description = "Terraform deployment at ${timestamp()}"
  depends_on = [module.connect, module.publish]
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_dynamodb_table" "connections" {
  name           = "${var.stage}_connections"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ID"
  stream_enabled = true
  stream_view_type = "NEW_AND_OLD_IMAGES"
  attribute {
    name = "ID"
    type = "S"
  }
  ttl {
    attribute_name = "LiveUntil"
    enabled        = true
  }
}

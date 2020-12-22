resource "aws_api_gateway_rest_api" "pooling" {
  name        = "${var.stage}_pooling"
}

resource "aws_api_gateway_stage" "pooling" {
  stage_name    = var.stage
  rest_api_id   = aws_api_gateway_rest_api.pooling.id
  deployment_id = aws_api_gateway_deployment.pooling.id
}

resource "aws_api_gateway_deployment" "pooling" {
  depends_on  = [module.reconnect]
  rest_api_id = aws_api_gateway_rest_api.pooling.id
  description = "Terraform deployment at ${timestamp()}"
  lifecycle {
    create_before_destroy = true
  }
}
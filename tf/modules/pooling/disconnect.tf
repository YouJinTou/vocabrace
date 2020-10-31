
resource "null_resource" "disconnect" {
  provisioner "local-exec" {
    command = "go build -o disconnect disconnect.go && build-lambda-zip -output ../../../../tf/modules/pooling/disconnect.zip disconnect ../config.${var.stage}.json"
    working_dir = "../../src/lambda/pooling/disconnect"
    environment = {
      GOOS = "linux"
      GOARCH = "amd64"
      CGO_ENABLED = "0"
    }
  }
}

module "disconnect" {
  source = "../lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "./disconnect.zip"
  function_name = "${var.stage}_disconnect"
  handler = "disconnect"
  environment_variables = {
    STAGE: var.stage
  }
  function_can_invoke_api_gateway = true
  api_gateway_can_invoke_function = true
  api_gateway_source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/$disconnect"
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
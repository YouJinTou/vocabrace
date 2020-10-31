resource "null_resource" "publish" {
  provisioner "local-exec" {
    command = "go build -o publish publish.go && build-lambda-zip -output ../../../../tf/modules/pooling/publish.zip publish ../config.${var.stage}.json"
    working_dir = "../../src/lambda/pooling/publish"
    environment = {
      GOOS = "linux"
      GOARCH = "amd64"
      CGO_ENABLED = "0"
    }
  }
}

module "publish" {
  source = "../lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "./publish.zip"
  function_name = "${var.stage}_publish"
  handler = "publish"
  environment_variables = {
    STAGE: var.stage
  }
  function_can_invoke_api_gateway = true
  api_gateway_can_invoke_function = true
  api_gateway_source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/publish"
  depends_on = [null_resource.publish]
}

resource "aws_apigatewayv2_route" "publish" {
  api_id    = aws_apigatewayv2_api.pooling.id
  route_key = "publish"
  target = "integrations/${aws_apigatewayv2_integration.publish.id}"
}

resource "aws_apigatewayv2_integration" "publish" {
  api_id           = aws_apigatewayv2_api.pooling.id
  integration_type = "AWS_PROXY"
  connection_type           = "INTERNET"
  integration_method        = "POST"
  integration_uri           = module.publish.this_lambda_function_invoke_arn
}
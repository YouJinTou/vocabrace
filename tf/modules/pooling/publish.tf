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
  source = "terraform-aws-modules/lambda/aws"
  function_name = "${var.stage}_publish"
  description   = "Invoked by the API Gateway Websocket runtime when a client publishes a message."
  handler       = "publish"
  runtime       = "go1.x"
  source_path = [
    {
      path = "${path.module}/publish.zip"
      pip_requirements = false
    }
  ]
  attach_policy = true
  policy = aws_iam_policy.pooling.arn
  create_current_version_allowed_triggers = false
  allowed_triggers = {
    APIGatewayPoolingConnect = {
      service = "apigateway"
      source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/$publish"
    }
  }
  environment_variables = {
    STAGE = var.stage
  }
  tags = {
    stage = var.stage
  }
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
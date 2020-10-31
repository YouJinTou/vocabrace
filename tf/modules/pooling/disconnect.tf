
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
  source = "terraform-aws-modules/lambda/aws"
  function_name = "${var.stage}_disconnect"
  description   = "Invoked by the API Gateway Websocket runtime when a client disconnects."
  handler       = "disconnect"
  runtime       = "go1.x"
  source_path = [
    {
      path = "${path.module}/disconnect.zip"
      pip_requirements = false
    }
  ]
  attach_policy = true
  policy = aws_iam_policy.pooling.arn
  create_current_version_allowed_triggers = false
  allowed_triggers = {
    APIGatewayPoolingConnect = {
      service = "apigateway"
      source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/$disconnect"
    }
  }
  environment_variables = {
    STAGE = var.stage
  }
  tags = {
    stage = var.stage
  }
  depends_on = [null_resource.disconnect]
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
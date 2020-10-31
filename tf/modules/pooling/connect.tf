
resource "null_resource" "connect" {
  provisioner "local-exec" {
    command = "go build -o connect connect.go && build-lambda-zip -output ../../../../tf/modules/pooling/connect.zip connect ../config.${var.stage}.json"
    working_dir = "../../src/lambda/pooling/connect"
    environment = {
      GOOS = "linux"
      GOARCH = "amd64"
      CGO_ENABLED = "0"
    }
  }
}

module "connect" {
  source = "terraform-aws-modules/lambda/aws"
  function_name = "${var.stage}_connect"
  description   = "Invoked by the API Gateway Websocket runtime when a client connects."
  handler       = "connect"
  runtime       = "go1.x"
  source_path = [
    {
      path = "${path.module}/connect.zip"
      pip_requirements = false
    }
  ]
  attach_cloudwatch_logs_policy = false
  attach_policy = true
  policy = aws_iam_policy.pooling.arn
  create_current_version_allowed_triggers = false
  allowed_triggers = {
    APIGatewayPoolingConnect = {
      service = "apigateway"
      source_arn = "${aws_apigatewayv2_api.pooling.execution_arn}/*/$connect"
    }
  }
  environment_variables = {
    STAGE = var.stage
  }
  tags = {
    stage = var.stage
  }
  depends_on = [null_resource.connect]
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
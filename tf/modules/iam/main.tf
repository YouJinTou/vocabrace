resource "aws_api_gateway_rest_api" "iam" {
  name        = "${var.stage}_iam"
}

resource "aws_api_gateway_stage" "iam" {
  stage_name    = var.stage
  rest_api_id   = aws_api_gateway_rest_api.iam.id
  deployment_id = aws_api_gateway_deployment.iam.id
}

resource "aws_api_gateway_deployment" "iam" {
  depends_on  = [module.iam]
  rest_api_id = aws_api_gateway_rest_api.iam.id
  description = "Terraform deployment at ${timestamp()}"
  lifecycle {
    create_before_destroy = true
  }
}

module "iam" {
  source = "../lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../payloads/iam.zip"
  function_name = "${var.stage}_iam"
  handler = "iam"
  timeout = 30
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
    ACCOUNT_ID: var.aws_account_id
  }
  api_gateway_can_invoke_function = true
  api_gateway_source_arn = "${aws_api_gateway_rest_api.iam.execution_arn}/*"
  rest_api_integration = {
        rest_api_id: aws_api_gateway_rest_api.iam.id,
        root_resource_id: aws_api_gateway_rest_api.iam.root_resource_id,
        path_parts: ["iam", "provider-auth"],
        http_methods: ["POST"]
    }
}
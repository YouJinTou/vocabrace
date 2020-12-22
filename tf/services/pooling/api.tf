module "get_user_pool" {
  source = "../../modules/lambda"
  aws_region = var.aws_region
  aws_account_id = var.aws_account_id
  filename = "../../payloads/getuserpool.zip"
  function_name = "${var.stage}_getuserpool"
  handler = "getuserpool"
  timeout = 30
  environment_variables = {
    STAGE: var.stage
    REGION: var.aws_region
    ACCOUNT_ID: var.aws_account_id
    IS_SERVERLESS: true
  }
  api_gateway_can_invoke_function = true
  api_gateway_source_arn = "${aws_api_gateway_rest_api.pooling.execution_arn}/*"
  rest_api_integration = {
        rest_api_id: aws_api_gateway_rest_api.pooling.id,
        root_resource_id: aws_api_gateway_rest_api.pooling.root_resource_id,
        path_parts: ["pooling", "userpools", "{userID}"],
        http_methods: ["GET"],
        enable_cors: true
    }
}
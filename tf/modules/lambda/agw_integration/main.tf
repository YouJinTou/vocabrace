locals {
  create1 = length(var.path_parts) > 0
  create2 = length(var.path_parts) > 1
  create3 = length(var.path_parts) > 2
  create4 = length(var.path_parts) > 3
  last_resource_id = (
    local.create4 ? aws_api_gateway_resource.resource4[0].id :
    local.create3 ? aws_api_gateway_resource.resource3[0].id :
    local.create2 ? aws_api_gateway_resource.resource2[0].id :
    local.create1 ? aws_api_gateway_resource.resource1[0].id :
    "")
}

resource "aws_api_gateway_resource" "resource1" {
  count = local.create1 ? 1 : 0
  path_part   = var.path_parts[0] 
  parent_id   = var.root_resource_id
  rest_api_id = var.rest_api_id
}

resource "aws_api_gateway_resource" "resource2" {
  count = local.create2 ? 1 : 0
  path_part   = var.path_parts[1] 
  parent_id   = aws_api_gateway_resource.resource1[0].id
  rest_api_id = var.rest_api_id
}

resource "aws_api_gateway_resource" "resource3" {
  count = local.create3 ? 1 : 0
  path_part   = var.path_parts[2] 
  parent_id   = aws_api_gateway_resource.resource2[0].id
  rest_api_id = var.rest_api_id
}

resource "aws_api_gateway_resource" "resource4" {
  count = local.create4 ? 1 : 0
  path_part   = var.path_parts[3] 
  parent_id   = aws_api_gateway_resource.resource3[0].id
  rest_api_id = var.rest_api_id
}

resource "aws_api_gateway_method" "method" {
  for_each = toset(var.http_methods)
  rest_api_id   = var.rest_api_id
  resource_id   = local.last_resource_id
  http_method   = each.key
  authorization = "NONE"
}

module "cors" {
  count = var.enable_cors ? 1 : 0
  source = "./cors"
  rest_api_id = var.rest_api_id
  resource_id = local.last_resource_id
}

resource "aws_api_gateway_integration" "integration" {
  for_each = { for idx, method in toset(var.http_methods): method => idx }
  rest_api_id             = var.rest_api_id
  resource_id             = local.last_resource_id
  http_method             = aws_api_gateway_method.method[each.value].http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.function_invoke_arn
}
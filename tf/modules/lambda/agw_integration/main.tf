resource "aws_api_gateway_resource" "resource" {
  for_each = toset(var.path_parts)
  path_part   = each.key
  parent_id   = local._4th_agw_parent
  rest_api_id = var.rest_api_id
}

resource "aws_api_gateway_method" "method" {
  for_each = toset(var.http_methods)
  rest_api_id   = var.rest_api_id
  resource_id   = local.last_agw_resource_id
  http_method   = each.key
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  for_each = { for idx, method in toset(var.http_methods): method => idx }
  rest_api_id             = var.rest_api_id
  resource_id             = local.last_agw_resource_id
  http_method             = aws_api_gateway_method.method[each.value].http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.function_invoke_arn
}
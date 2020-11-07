resource "aws_iam_role" "role" {
  name = var.function_name
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "cloudwatch" {
  name = "${var.function_name}_cloudwatch"
  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
          "Effect": "Allow",
          "Action": [
              "logs:CreateLogStream",
              "logs:PutLogEvents"
          ],
          "Resource": "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:/aws/lambda/${aws_lambda_function.function.function_name}:*"
      },
      {
          "Effect": "Allow",
          "Action": [
              "logs:CreateLogGroup"
          ],
          "Resource": "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:*"
      }
    ]
  }
  EOF
}

resource "aws_iam_role_policy_attachment" "cloudwatch" {
  role       = aws_iam_role.role.name
  policy_arn = aws_iam_policy.cloudwatch.arn
}

resource "aws_iam_role_policy_attachment" "api_gateway" {
  count = var.function_can_invoke_api_gateway ? 1 : 0
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess"
}

resource "aws_iam_role_policy_attachment" "administrator" {
  count = var.is_administrator ? 1 : 0
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

resource "aws_lambda_function" "function" {
  filename      = var.filename
  function_name = var.function_name
  role          = aws_iam_role.role.arn
  handler       = var.handler
  source_code_hash = filebase64sha256(var.filename)
  runtime = "go1.x"
  timeout = var.timeout
  reserved_concurrent_executions = var.reserved_concurrent_executions
  environment {
    variables = var.environment_variables
  }
}

resource "aws_lambda_permission" "api_gateway_permission" {
  count = var.api_gateway_can_invoke_function ? 1 : 0
  statement_id_prefix  = "${var.function_name}_AllowInvocationsFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = var.api_gateway_source_arn
}

resource "aws_lambda_event_source_mapping" "sqs" {
  count = var.sqs_can_invoke_function ? 1 : 0
  batch_size        = 10
  event_source_arn  = var.sqs_source_arn
  function_name     = aws_lambda_function.function.function_name
}

resource "aws_iam_role_policy_attachment" "sqs" {
  count = var.sqs_can_invoke_function ? 1 : 0
  role       = aws_iam_role.role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_cloudwatch_event_target" "target" {
  count = var.cloudwatch_can_invoke_function ? 1 : 0
  rule      = aws_cloudwatch_event_rule.rule[0].name
  arn       = aws_lambda_function.function.arn
}

resource "aws_cloudwatch_event_rule" "rule" {
  count = var.cloudwatch_can_invoke_function ? 1 : 0
  name_prefix        = aws_lambda_function.function.function_name
  schedule_expression = var.cloudwatch_event_rule
}

resource "aws_lambda_permission" "rule_permission" {
  count = var.cloudwatch_can_invoke_function ? 1 : 0
  statement_id  = "AllowExecutionFromCloudWatchRule"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.rule[0].arn
}
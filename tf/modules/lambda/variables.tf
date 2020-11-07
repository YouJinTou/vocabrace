variable aws_account_id {
    type = string
}

variable aws_region {
    type = string
}

variable is_administrator {
    type = bool
    default = true
}

variable filename {
    type = string
}

variable function_name {
    type = string
}

variable environment_variables {
    type = map
    default = {}
}

variable handler {
    type = string
}

variable function_can_invoke_api_gateway {
    type = bool
    default = false
}

variable api_gateway_can_invoke_function {
    type = bool
    default = false
}

variable api_gateway_source_arn {
    type = string
    default = ""
}

variable sqs_can_invoke_function {
    type = bool
    default = false
}

variable sqs_source_arn {
    type = string
    default = ""
}

variable dynamodb_can_invoke_function {
    type = bool
    default = false
}

variable dynamodb_stream_arn {
    type = string
    default  = ""
}
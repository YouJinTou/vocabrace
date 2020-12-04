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

variable rest_api_integration {
    type = object({
        rest_api_id: string,
        root_resource_id: string,
        path_parts: list(string),
        http_methods: list(string),
        enable_cors: bool
    })
    default = {
        rest_api_id: "",
        root_resource_id: "",
        path_parts: [],
        http_methods: [],
        enable_cors: false
    }
}

variable api_gateway_source_arn {
    type = string
    default = ""
}

variable sqs_sources {
    type = list(object({
        arn = string,
        batch_size = number
    }))
    default = []
}

variable cloudwatch_can_invoke_function {
    type = bool
    default = false
}

variable cloudwatch_event_rule {
    type = string
    default = ""
}

variable timeout {
    type = number
    default = 3
}

variable reserved_concurrent_executions {
    type = number
    default = -1
}
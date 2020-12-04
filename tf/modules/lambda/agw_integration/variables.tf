variable rest_api_id {
    type = string
}

variable root_resource_id {
    type = string
}

variable path_parts {
    type = list(string)
}

variable enable_cors {
    type = bool
    default = false
}

variable http_methods {
    type = list(string)
}

variable function_invoke_arn {
    type = string
}
variable stage {
    type = string
    description = "The staging environment."
}

variable aws_account_id {
    type = string
}

variable aws_region {
    type = string
}

variable conductor_queue_arn {
    type = string
    description = "The conductor's queue ARN."
}
output topic_arn {
    value = aws_sns_topic.facts.arn
}

output table_arn {
    value = aws_dynamodb_table.facts.arn
}
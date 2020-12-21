output topic_arn {
    value = aws_sns_topic.events.arn
}

output table_arn {
    value = aws_dynamodb_table.store.arn
}
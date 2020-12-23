
resource "aws_dynamodb_table" "wordlines_bulgarian" {
  name           = "wordlines_bulgarian"
  read_capacity  = 5
  write_capacity = 1
  hash_key       = "Word"
  attribute {
    name = "Word"
    type = "S"
  }
}

resource "aws_dynamodb_table" "wordlines_english" {
  name           = "wordlines_english"
  read_capacity  = 5
  write_capacity = 1
  hash_key       = "Word"
  attribute {
    name = "Word"
    type = "S"
  }
}

resource "aws_dynamodb_table" "missing_words" {
  name           = "wordlines_missing_words"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Word"
  attribute {
    name = "Word"
    type = "S"
  }
  ttl {
    attribute_name = "LiveUntil"
    enabled        = true
  }
}
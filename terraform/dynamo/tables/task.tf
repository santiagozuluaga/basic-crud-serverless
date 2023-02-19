resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "task"
  hash_key       = "id"
  range_key      = "title"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "title"
    type = "S"
  }
}

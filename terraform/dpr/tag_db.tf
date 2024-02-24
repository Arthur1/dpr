resource "aws_dynamodb_table" "tag_db" {
  name           = var.tag_db_dynamodb_table_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 3
  write_capacity = 3

  attribute {
    name = "tag"
    type = "S"
  }
  attribute {
    name = "object_key"
    type = "S"
  }
  hash_key = "tag"

  global_secondary_index {
    name            = "index_object_key_tag"
    hash_key        = "object_key"
    range_key       = "tag"
    projection_type = "ALL"
    read_capacity   = 3
    write_capacity  = 3
  }

  tags = var.tags
}

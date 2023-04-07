resource "harperdb_schema" "dogs" {
  name = "dogs2"
}

resource "harperdb_table" "table_1" {
  schema = harperdb_schema.dogs.name
  name = "table_1"
  hash_attribute = "attributed"
}
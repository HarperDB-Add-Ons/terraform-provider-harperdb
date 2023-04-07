resource "harperdb_role" "basic" {
  name = "reader_role"
  schema_permissions = {
    dogs2 = {
      tables = {
        "${harperdb_table.table_1.name}" = {
          read = true
          delete = true
          attribute_permissions = [
            {
              name = "attributed"
              read = false
            }
          ]
        }
      }
    }
  }

  # depends_on = [ harperdb_table.table_1 ]
}
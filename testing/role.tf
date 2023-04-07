resource "harperdb_role" "basic" {
  name = "reader_role"
  super_user = false
  cluster_user = false
  schema_permissions = {
    dogs2 = {
      tables = {
        "${harperdb_table.table_1.name}" = {
          read = true
          delete = true
          attribute_permissions = [
            {
              name = "attributed"
              read = true
            }
          ]
        }
      }
    }
  }

  # depends_on = [ harperdb_table.table_1 ]
}
resource "harperdb_permission" "basic" {
    super_user = true
    cluster_user = false
    table_permissions = {
        testing = {}
        invalid = {}
    }
}
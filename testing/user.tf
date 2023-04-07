resource "harperdb_user" "hamish" {
    username = "hamish"
    password = "password"
    active = true
    role = harperdb_role.basic.name
}
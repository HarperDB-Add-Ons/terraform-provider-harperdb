package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRole2Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "name",
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig("Role1", "password", "", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_Role.test", "name", "Role1"),
					resource.TestCheckResourceAttr("harperdb_Role.test", "role", ""),
				),
			},
			{
				Config: testAccRoleResourceConfig("Role2", "password", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_Role.test", "name", "two"),
				),
			},
		},
	})
}

func testAccRole2ResourceConfig(Role, pass, role string, active bool) string {
	return fmt.Sprintf(`
	%s

resource "harperdb_role" "test" {
  Rolename = "%s"
  password = "%s"
  role = "%s"
  active = %t
}
`, testAccProviderTF(), Role, pass, role, active)
}

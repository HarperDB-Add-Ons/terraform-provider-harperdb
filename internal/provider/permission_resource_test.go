package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPermissionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "name",
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccPermissionResourceConfig("Permission1", "password", "", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_permission.test", "name", "Permission1"),
					resource.TestCheckResourceAttr("harperdb_permission.test", "Permission", ""),
				),
			},
			{
				Config: testAccPermissionResourceConfig("Permission2", "password", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_Permission.test", "name", "two"),
				),
			},
		},
	})
}

func testAccPermissionResourceConfig(user, pass, role string, active bool) string {
	return fmt.Sprintf(`
	%s

resource "harperdb_Permission" "test" {
  Permissionname = "%s"
  password = "%s"
  Permission = "%s"
  active = %t
}
`, testAccProviderTF(), user, pass, role, active)
}

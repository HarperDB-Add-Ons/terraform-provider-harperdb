package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "name",
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig("user1", "password", "", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_user.test", "name", "user1"),
					resource.TestCheckResourceAttr("harperdb_user.test", "role", ""),
				),
			},
			{
				Config: testAccUserResourceConfig("user2", "password", "", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_user.test", "name", "two"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(user, pass, role string, active bool) string {
	return fmt.Sprintf(`
	%s

resource "harperdb_user" "test" {
  username = "%s"
  password = "%s"
  role = "%s"
  active = %t
}
`, testAccProviderTF(), user, pass, role, active)
}

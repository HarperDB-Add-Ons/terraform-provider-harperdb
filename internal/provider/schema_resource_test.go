package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "name",
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccSchemaResourceConfig("https://terraform-test-moredhel.harperdbcloud.com", "moredhel", "question", "one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_schema.test", "name", "one"),
					resource.TestCheckResourceAttr("harperdb_schema.test", "id", "one"),
				),
			},
			{
				Config: testAccSchemaResourceConfig("https://terraform-test-moredhel.harperdbcloud.com", "moredhel", "question", "two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("harperdb_schema.test", "name", "two"),
				),
			},
		},
	})
}

func testAccProviderTF(host, user, pass string) string {
	return fmt.Sprintf(`
		provider "harperdb" {
			endpoint = "%s"
			username = "%s"
			password = "%s"
		}
	`, host, user, pass)
}

func testAccSchemaResourceConfig(host, user, pass, configurableAttribute string) string {
	return fmt.Sprintf(`
	%s

resource "harperdb_schema" "test" {
  name = "%s"
}
`, testAccProviderTF(host, user, pass), configurableAttribute)
}

package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSandboxPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        gitlabProvider,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + gitlabTestingDefinition + `
resource "crczp_sandbox_pool" "test" {
  definition = {
    id = crczp_sandbox_definition.test.id
  }
  max_size = 2
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "size", "0"),
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "max_size", "2"),
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "rev", os.Getenv("TF_VAR_TAG_NAME")),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "rev_sha"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.vcpu"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.ram"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.instances"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.network"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.subnet"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.port"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.sub"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.full_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.given_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.family_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.mail"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.id",
						"crczp_sandbox_pool.test", "definition.created_by.id"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.sub",
						"crczp_sandbox_pool.test", "definition.created_by.sub"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.full_name",
						"crczp_sandbox_pool.test", "definition.created_by.full_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.given_name",
						"crczp_sandbox_pool.test", "definition.created_by.given_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.family_name",
						"crczp_sandbox_pool.test", "definition.created_by.family_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.mail",
						"crczp_sandbox_pool.test", "definition.created_by.mail"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "crczp_sandbox_pool.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + gitlabTestingDefinition + `
resource "crczp_sandbox_pool" "test" {
  definition = {
    id = crczp_sandbox_definition.test.id
  }
  max_size = 10
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "size", "0"),
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "max_size", "10"),
					resource.TestCheckResourceAttr("crczp_sandbox_pool.test", "rev", os.Getenv("TF_VAR_TAG_NAME")),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "rev_sha"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.vcpu"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.ram"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.instances"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.network"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.subnet"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "hardware_usage.port"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.sub"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.full_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.given_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.family_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_pool.test", "created_by.mail"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.id",
						"crczp_sandbox_pool.test", "definition.created_by.id"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.sub",
						"crczp_sandbox_pool.test", "definition.created_by.sub"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.full_name",
						"crczp_sandbox_pool.test", "definition.created_by.full_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.given_name",
						"crczp_sandbox_pool.test", "definition.created_by.given_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.family_name",
						"crczp_sandbox_pool.test", "definition.created_by.family_name"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_pool.test", "created_by.mail",
						"crczp_sandbox_pool.test", "definition.created_by.mail"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

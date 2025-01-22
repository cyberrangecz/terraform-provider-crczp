package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSandboxAllocationUnitResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + gitlabTestingDefinition + `
resource "crczp_sandbox_pool" "test" {
  definition = {
    id = crczp_sandbox_definition.test.id
  }
  max_size = 1
}

resource "crczp_sandbox_allocation_unit" "test" {
  pool_id = crczp_sandbox_pool.test.id
}

data "crczp_sandbox_request_output" "test-user" {
  id = crczp_sandbox_allocation_unit.test.allocation_request.id
}
data "crczp_sandbox_request_output" "test-networking" {
  id = crczp_sandbox_allocation_unit.test.allocation_request.id
  stage = "networking-ansible"
}
data "crczp_sandbox_request_output" "test-terraform" {
  id = crczp_sandbox_allocation_unit.test.allocation_request.id
  stage = "terraform"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_sandbox_allocation_unit.test", "locked", "false"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "id"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_allocation_unit.test", "pool_id",
						"crczp_sandbox_pool.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "allocation_request.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "allocation_request.allocation_unit_id"),
					resource.TestCheckResourceAttrPair("crczp_sandbox_allocation_unit.test", "allocation_request.allocation_unit_id",
						"crczp_sandbox_allocation_unit.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "allocation_request.created"),
					resource.TestCheckResourceAttr("crczp_sandbox_allocation_unit.test", "allocation_request.stages.#", "3"),
					resource.TestCheckResourceAttr("crczp_sandbox_allocation_unit.test", "allocation_request.stages.0", "FINISHED"),
					resource.TestCheckResourceAttr("crczp_sandbox_allocation_unit.test", "allocation_request.stages.1", "FINISHED"),
					resource.TestCheckResourceAttr("crczp_sandbox_allocation_unit.test", "allocation_request.stages.2", "FINISHED"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.sub"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.full_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.given_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.family_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_allocation_unit.test", "created_by.mail"),

					// Datasource sandbox request output
					resource.TestCheckResourceAttrPair("data.crczp_sandbox_request_output.test-user", "id",
						"crczp_sandbox_allocation_unit.test", "allocation_request.id"),
					resource.TestCheckResourceAttr("data.crczp_sandbox_request_output.test-user", "stage", "user-ansible"),
					resource.TestCheckResourceAttrSet("data.crczp_sandbox_request_output.test-user", "result"),

					resource.TestCheckResourceAttrPair("data.crczp_sandbox_request_output.test-networking", "id",
						"crczp_sandbox_allocation_unit.test", "allocation_request.id"),
					resource.TestCheckResourceAttr("data.crczp_sandbox_request_output.test-networking", "stage", "networking-ansible"),
					resource.TestCheckResourceAttrSet("data.crczp_sandbox_request_output.test-networking", "result"),

					resource.TestCheckResourceAttrPair("data.crczp_sandbox_request_output.test-terraform", "id",
						"crczp_sandbox_allocation_unit.test", "allocation_request.id"),
					resource.TestCheckResourceAttr("data.crczp_sandbox_request_output.test-terraform", "stage", "terraform"),
					resource.TestCheckResourceAttrSet("data.crczp_sandbox_request_output.test-terraform", "result"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "crczp_sandbox_allocation_unit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

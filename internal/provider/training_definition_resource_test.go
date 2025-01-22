package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const ltdDefinition = `
{
 "title" : "test",
 "description" : null,
 "prerequisites" : [ ],
 "outcomes" : [ ],
 "state" : "UNRELEASED",
 "levels" : [ ],
 "estimated_duration" : 0,
 "variant_sandboxes" : false
}
`

func TestAccTrainingDefinitionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "crczp_training_definition" "test" {
 content = <<EOL
` + ltdDefinition + `EOL
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_training_definition.test", "content", ltdDefinition),
					resource.TestCheckResourceAttrSet("crczp_training_definition.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "crczp_training_definition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

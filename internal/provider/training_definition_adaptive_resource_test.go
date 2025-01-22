package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const atdDefinition = `
{
 "title" : "test",
 "description" : null,
 "prerequisites" : [ ],
 "outcomes" : [ ],
 "state" : "UNRELEASED",
 "phases" : [ ],
 "estimated_duration" : 0
}
`

func TestAccTrainingDefinitionAdaptiveResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "crczp_training_definition_adaptive" "test" {
 content = <<EOL
` + atdDefinition + `EOL
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_training_definition_adaptive.test", "content", atdDefinition),
					resource.TestCheckResourceAttrSet("crczp_training_definition_adaptive.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "crczp_training_definition_adaptive.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

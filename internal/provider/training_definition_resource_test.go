package provider_test

//
//import (
//	"strconv"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//)
//
//const ltdDefinition = "{\"description\":null,\"estimated_duration\":0,\"levels\":[],\"outcomes\":[],\"prerequisites\":[]," +
//	"\"state\":\"UNRELEASED\",\"title\":\"test\",\"variant_sandboxes\":false}"
//
//func TestAccTrainingDefinitionResource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: providerConfig + `
//resource "crczp_training_definition" "test" {
// content = ` + strconv.Quote(ltdDefinition) + `
//}`,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("crczp_training_definition.test", "content", ltdDefinition),
//					resource.TestCheckResourceAttrSet("crczp_training_definition.test", "id"),
//				),
//			},
//			// ImportState testing
//			{
//				ResourceName:      "crczp_training_definition.test",
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}

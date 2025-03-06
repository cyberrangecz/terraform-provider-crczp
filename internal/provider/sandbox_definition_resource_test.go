package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testingDefinitionTag = `
variable "TAG_NAME" {}
variable "TOKEN" {}

resource "terraform_data" "git_tag" {
  count = 2
  input = {
    tag_name = var.TAG_NAME
  }
  provisioner "local-exec" {
    command = <<EOT
    GIT_SSH_COMMAND='ssh -o IdentitiesOnly=yes' git clone https://${var.TOKEN}@github.com/cyberrangecz/terraform-testing-definition.git repo
    cd repo
    git tag ${self.input.tag_name}-${count.index} -m ""
    git push origin ${self.input.tag_name}-${count.index}
    EOT
  }
  provisioner "local-exec" {
    when    = destroy
    command = <<EOT
    cd repo
    git push --delete origin ${self.input.tag_name}-${count.index}
    EOT
  }
}
`

func TestAccSandboxDefinitionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testingDefinitionTag + `
resource "crczp_sandbox_definition" "test" {
  url = "https://github.com/cyberrangecz/terraform-testing-definition.git"
  rev = "${terraform_data.git_tag[0].output.tag_name}-0"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "url", "https://github.com/cyberrangecz/terraform-testing-definition.git"),
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "rev", os.Getenv("TF_VAR_TAG_NAME")+"-0"),
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "name", "terraform-testing-definition"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.sub"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.full_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.given_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.family_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.mail"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "crczp_sandbox_definition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testingDefinitionTag + `
resource "crczp_sandbox_definition" "test" {
  url = "https://github.com/cyberrangecz/terraform-testing-definition.git"
  rev = "${terraform_data.git_tag[1].output.tag_name}-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "url", "https://github.com/cyberrangecz/terraform-testing-definition.git"),
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "rev", os.Getenv("TF_VAR_TAG_NAME")+"-1"),
					resource.TestCheckResourceAttr("crczp_sandbox_definition.test", "name", "terraform-testing-definition"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.id"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.sub"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.full_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.given_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.family_name"),
					resource.TestCheckResourceAttrSet("crczp_sandbox_definition.test", "created_by.mail"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

package provider_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"terraform-provider-crczp/internal/provider"
)

const (
	providerConfig = `
provider "crczp" {
  endpoint    = "https://stage.crp.kypo.muni.cz"
  retry_count = 3
}
`

	testingDefinition = `
variable "TAG_NAME" {}
variable "TOKEN" {}

resource "terraform_data" "git_tag" {
  input = {
    tag_name = var.TAG_NAME
  }
  provisioner "local-exec" {
    command = <<EOT
    GIT_SSH_COMMAND='ssh -o IdentitiesOnly=yes' git clone https://${var.TOKEN}@github.com/cyberrangecz/terraform-testing-definition.git repo
    cd repo
    git tag ${self.input.tag_name} -m ""
    git push origin ${self.input.tag_name}
    EOT
  }
  provisioner "local-exec" {
    when    = destroy
    command = <<EOT
    cd repo
    git push --delete origin ${self.input.tag_name}
    EOT
  }
}

resource "crczp_sandbox_definition" "test" {
  url = "https://github.com/cyberrangecz/terraform-testing-definition.git"
  rev = terraform_data.git_tag.output.tag_name
}
`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"crczp": providerserver.NewProtocol6WithError(provider.New("test")()),
}

//func testAccPreCheck(t *testing.T) {
// You can add code here to run prior to any test case execution, for example assertions
// about the appropriate environment variables being set are common to see in a pre-check
// function.
//}

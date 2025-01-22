package provider_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"terraform-provider-crczp/internal/provider"
)

const (
	providerConfig = `
provider "crczp" {
  endpoint    = "https://lab.crp.crczp.muni.cz"
  retry_count = 3
}
`

	gitlabTestingDefinition = `
variable "TAG_NAME" {}

resource "null_resource" "git_tag" {
  provisioner "local-exec" {
    command = <<EOT
    cat <<EOF > key
    -----BEGIN OPENSSH PRIVATE KEY-----
    b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
    QyNTUxOQAAACBcy6PgN52E5RRdEvPIyrRzWGGB00z0htPTZfTZHLSdjAAAAJg9eMq9PXjK
    vQAAAAtzc2gtZWQyNTUxOQAAACBcy6PgN52E5RRdEvPIyrRzWGGB00z0htPTZfTZHLSdjA
    AAAEBD2CRf7TB/rCgGdryTHa3S0bg0Z2QE/tshWZEi+Izzg1zLo+A3nYTlFF0S88jKtHNY
    YYHTTPSG09Nl9NkctJ2MAAAAD3pkZW5la0Bza2VsbGlnZQECAwQFBg==
    -----END OPENSSH PRIVATE KEY-----
    EOF
    chmod 600 key
    GIT_SSH_COMMAND='ssh -i key -o IdentitiesOnly=yes' git clone git@github.com:cyberrangecz/terraform-testing-definition.git repo
    cd repo
    git tag ${var.TAG_NAME}
    git push origin ${var.TAG_NAME}
    EOT
  }
  provisioner "local-exec" {
    when    = destroy
    command = <<EOT
    cat <<EOF > key
    -----BEGIN OPENSSH PRIVATE KEY-----
    b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
    QyNTUxOQAAACBcy6PgN52E5RRdEvPIyrRzWGGB00z0htPTZfTZHLSdjAAAAJg9eMq9PXjK
    vQAAAAtzc2gtZWQyNTUxOQAAACBcy6PgN52E5RRdEvPIyrRzWGGB00z0htPTZfTZHLSdjA
    AAAEBD2CRf7TB/rCgGdryTHa3S0bg0Z2QE/tshWZEi+Izzg1zLo+A3nYTlFF0S88jKtHNY
    YYHTTPSG09Nl9NkctJ2MAAAAD3pkZW5la0Bza2VsbGlnZQECAwQFBg==
    -----END OPENSSH PRIVATE KEY-----
    EOF
    chmod 600 key
    GIT_SSH_COMMAND='ssh -i key -o IdentitiesOnly=yes' git clone git@github.com:cyberrangecz/terraform-testing-definition.git repo || true
    cd repo
    git push --delete origin ${var.TAG_NAME}
    EOT
  }
}

resource "crczp_sandbox_definition" "test" {
  url = "https://gitlab.ics.muni.cz/muni-crczp-crp/prototypes-and-examples/sandbox-definitions/terraform-provider-testing-definition.git"
  rev = TAG_NAME
  depends_on = [
    null_resource.git_tag
  ]
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

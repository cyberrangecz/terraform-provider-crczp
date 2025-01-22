# Terraform CRCZP Provider

Terraform provider for CRCZP allows [Terraform](https://www.terraform.io/) to manage [CRCZP]() resources.

See documentation at the [Terraform Registry]().

## Example Usage

```terraform
terraform {
  required_providers {
    crczp = {
      source = "cyberrangecz/crczp"
      version = "0.3.1"
    }
  }
}

provider "crczp" {
  endpoint  = "https://your.crczp.ex" # Or use CRCZP_ENDPOINT env var
  client_id = "***"                  # Or use CRCZP_CLIENT_ID env var
  username  = "user"                 # Or use CRCZP_USERNAME env var
  password  = "***"                  # Or use CRCZP_PASSWORD env var
}

resource "crczp_sandbox_definition" "example" {
  url = "https://github.com/cyberrangecz/library-junior-hacker.git"
  rev = "master"
}

resource "crczp_sandbox_pool" "example" {
  definition = {
    id = crczp_sandbox_definition.example.id
  }
  max_size = 1
}

resource "crczp_sandbox_allocation_unit" "example" {
  pool_id = crczp_sandbox_pool.example.id
}
```

# License
Parts of the Source Code, namely [GNUmakefile](./GNUmakefile), [.golangci.yml](./.golangci.yml), [.goreleaser.yml](./.goreleaser.yml),
[tools/tools.go](./tools/tools.go),
[internal/validators/timeduration.go](./internal/validators/timeduration.go) and [internal/validators/timeduration_test.go](./internal/validators/timeduration_test.go),
is subject to the terms of the Mozilla Public License, v. 2.0. You can obtain the license at https://mozilla.org/MPL/2.0/. 
The Source Code was used from [HashiCorp terraform-provider-scaffolding-framework](https://github.com/hashicorp/terraform-provider-scaffolding-framework) and [HashiCorp terraform-plugin-framework-timeouts](https://github.com/hashicorp/terraform-plugin-framework-timeouts) repositories.

The rest of the Source Code is subject to the [MIT License](./LICENSE).


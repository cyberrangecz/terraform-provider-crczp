resource "crczp_sandbox_definition" "example" {
  url = "git@gitlab.ics.muni.cz:muni-crczp-trainings/games/junior-hacker.git"
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

data "crczp_sandbox_request_output" "example" {
  id = crczp_sandbox_allocation_unit.example.allocation_request.id
}

output "example" {
  value = data.crczp_sandbox_request_output.example.result
}

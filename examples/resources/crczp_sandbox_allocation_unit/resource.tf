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

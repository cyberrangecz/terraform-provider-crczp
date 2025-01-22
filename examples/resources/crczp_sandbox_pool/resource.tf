resource "crczp_sandbox_definition" "example" {
  url = "https://github.com/cyberrangecz/library-junior-hacker.git"
  rev = "master"
}

resource "crczp_sandbox_pool" "example" {
  definition = {
    id = crczp_sandbox_definition.example.id
  }
  max_size = 2
}

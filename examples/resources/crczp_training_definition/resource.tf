resource "crczp_training_definition" "example" {
  content = jsonencode(
    {
      title              = "test"
      description        = null
      state              = "UNRELEASED"
      variant_sandboxes  = false
      estimated_duration = 0
      levels             = []
      outcomes           = []
      prerequisites      = []
    }
  )
}

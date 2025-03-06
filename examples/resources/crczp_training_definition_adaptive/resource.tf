resource "crczp_training_definition_adaptive" "example" {
  content = jsonencode(
    {
      title              = "test"
      description        = null
      state              = "UNRELEASED"
      estimated_duration = 0
      outcomes           = []
      phases             = []
      prerequisites      = []
    }
  )
}

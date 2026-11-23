package provider_test

import (
	"strconv"
	"testing"

	"terraform-provider-crczp/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ltdDefinition is in the shape the backend currently exports.
const ltdDefinition = "\n{\"title\":\"test\",\"description\":\"test\",\"prerequisites\":[],\"outcomes\":[],\"levels\":[]}\n"

// legacyDefinition is in the shape older backends exported. It carries the state, estimated_duration,
// show_stepper_bar, variant_sandboxes and per-level order properties the current export leaves out.
const legacyDefinition = `{
  "title": "test-legacy",
  "description": "Definition exported by an older backend",
  "prerequisites": [],
  "outcomes": ["outcome"],
  "state": "UNRELEASED",
  "show_stepper_bar": true,
  "levels": [
    {
      "title": "Introduction",
      "level_type": "INFO_LEVEL",
      "order": 0,
      "estimated_duration": 0,
      "minimal_possible_solve_time": null,
      "content": "First level"
    },
    {
      "title": "Conclusion",
      "level_type": "INFO_LEVEL",
      "order": 1,
      "estimated_duration": 1,
      "minimal_possible_solve_time": null,
      "content": "Second level"
    }
  ],
  "estimated_duration": 1,
  "variant_sandboxes": false
}`

// legacyDefinitionExport is legacyDefinition as the current backend exports it.
const legacyDefinitionExport = `{"title":"test-legacy","description":"Definition exported by an older backend",` +
	`"prerequisites":[],"outcomes":["outcome"],"levels":[` +
	`{"title":"Introduction","level_type":"INFO_LEVEL","estimated_duration":0,"minimal_possible_solve_time":null,"content":"First level"},` +
	`{"title":"Conclusion","level_type":"INFO_LEVEL","estimated_duration":1,"minimal_possible_solve_time":null,"content":"Second level"}]}`

func TestAccTrainingDefinitionResource(t *testing.T) {
	for name, content := range map[string]string{
		"current": ltdDefinition,
		"legacy":  legacyDefinition,
	} {
		t.Run(name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Create and Read testing
					{
						Config: providerConfig + `
resource "crczp_training_definition" "test" {
 content = ` + strconv.Quote(content) + `
}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("crczp_training_definition.test", "content", content),
							resource.TestCheckResourceAttrSet("crczp_training_definition.test", "id"),
						),
					},
					// ImportState testing
					{
						ResourceName: "crczp_training_definition.test",
						ImportState:  true,
					},
					// Delete testing automatically occurs in TestCase
				},
			})
		})
	}
}

func TestExportMatchesContent(t *testing.T) {
	tests := map[string]struct {
		exported string
		stored   string
		want     bool
	}{
		"identical": {
			exported: ltdDefinition,
			stored:   ltdDefinition,
			want:     true,
		},
		"formatting and key order differ": {
			exported: `{"levels":[],"outcomes":[],"prerequisites":[],"description":"test","title":"test"}`,
			stored:   ltdDefinition,
			want:     true,
		},
		"legacy content with properties no longer exported": {
			exported: legacyDefinitionExport,
			stored:   legacyDefinition,
			want:     true,
		},
		"level changed outside Terraform": {
			exported: `{"title":"test-legacy","description":"Definition exported by an older backend",` +
				`"prerequisites":[],"outcomes":["outcome"],"levels":[` +
				`{"title":"Introduction","level_type":"INFO_LEVEL","estimated_duration":0,"minimal_possible_solve_time":null,"content":"Changed"},` +
				`{"title":"Conclusion","level_type":"INFO_LEVEL","estimated_duration":1,"minimal_possible_solve_time":null,"content":"Second level"}]}`,
			stored: legacyDefinition,
			want:   false,
		},
		"level removed outside Terraform": {
			exported: `{"title":"test-legacy","description":"Definition exported by an older backend",` +
				`"prerequisites":[],"outcomes":["outcome"],"levels":[` +
				`{"title":"Introduction","level_type":"INFO_LEVEL","estimated_duration":0,"minimal_possible_solve_time":null,"content":"First level"}]}`,
			stored: legacyDefinition,
			want:   false,
		},
		"levels reordered outside Terraform": {
			exported: `{"title":"test-legacy","description":"Definition exported by an older backend",` +
				`"prerequisites":[],"outcomes":["outcome"],"levels":[` +
				`{"title":"Conclusion","level_type":"INFO_LEVEL","estimated_duration":1,"minimal_possible_solve_time":null,"content":"Second level"},` +
				`{"title":"Introduction","level_type":"INFO_LEVEL","estimated_duration":0,"minimal_possible_solve_time":null,"content":"First level"}]}`,
			stored: legacyDefinition,
			want:   false,
		},
		"exported property missing from stored content": {
			exported: `{"title":"test","description":null,"prerequisites":[],"outcomes":[],"levels":[]}`,
			stored:   `{"title":"test","prerequisites":[],"outcomes":[],"levels":[]}`,
			want:     false,
		},
		"type differs": {
			exported: `{"title":"test","outcomes":[]}`,
			stored:   `{"title":"test","outcomes":{}}`,
			want:     false,
		},
		"invalid stored content": {
			exported: ltdDefinition,
			stored:   `{`,
			want:     false,
		},
		"invalid export": {
			exported: `{`,
			stored:   ltdDefinition,
			want:     false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := provider.ExportMatchesContent(tt.exported, tt.stored); got != tt.want {
				t.Errorf("ExportMatchesContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

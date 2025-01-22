package plan_modifiers

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JSONNormalizePlanModifier struct{}

var _ planmodifier.String = JSONNormalizePlanModifier{}

func (m JSONNormalizePlanModifier) Description(_ context.Context) string {
	return "Normalizes the JSON string to remove whitespace and ensure consistent formatting."
}

func (m JSONNormalizePlanModifier) MarkdownDescription(_ context.Context) string {
	return "Normalizes the JSON string to remove whitespace and ensure consistent formatting."
}

func (m JSONNormalizePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	jsonInput := req.PlanValue.ValueString()

	normalized, diags := NormalizeJSON(jsonInput)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = types.StringValue(normalized)
}

func NormalizeJSON(input string) (string, diag.Diagnostics) {
	if input == "" {
		return input, nil
	}

	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(input), &parsed)
	if err != nil {
		return "", diag.Diagnostics{diag.NewErrorDiagnostic(
			"Invalid JSON", "Unable to parse JSON: "+err.Error(),
		)}
	}

	normalized, err := json.Marshal(parsed)
	if err != nil {
		return "", diag.Diagnostics{diag.NewErrorDiagnostic(
			"JSON Normalization Error", "Unable to normalize JSON: "+err.Error(),
		)}
	}

	return string(normalized), nil
}

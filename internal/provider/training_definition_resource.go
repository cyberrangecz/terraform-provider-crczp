package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cyberrangecz/go-client/pkg/crczp"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &trainingDefinitionResource{}
var _ resource.ResourceWithImportState = &trainingDefinitionResource{}
var _ resource.ResourceWithConfigure = &trainingDefinitionResource{}

func NewTrainingDefinitionResource() resource.Resource {
	return &trainingDefinitionResource{}
}

// trainingDefinitionResource defines the resource implementation.
type trainingDefinitionResource struct {
	client *crczp.Client
}

func (r *trainingDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_training_definition"
}

func (r *trainingDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Training definition",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Id of the training definition",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "JSON with exported training definition. Properties the backend no longer exports, such as `state`, `estimated_duration`, `variant_sandboxes` or a level's `order`, are ignored when detecting changes made outside Terraform.",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *trainingDefinitionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*crczp.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected crczp.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *trainingDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var content string

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("content"), &content)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition, err := r.client.CreateTrainingDefinition(ctx, content)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create training definition, got error: %s", err))
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created training definition %d", definition.Id))

	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *trainingDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.Int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	definition, err := r.client.GetTrainingDefinition(ctx, id.ValueInt64())
	if errors.Is(err, crczp.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read training definition, got error: %s", err))
		return
	}

	var storedContent jsontypes.Normalized

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("content"), &storedContent)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Newer backends export fewer properties than older exports carry, so keep the stored content
	// unless the export disagrees with it on a property the export still has.
	if !storedContent.IsNull() && !storedContent.IsUnknown() &&
		exportMatchesContent(definition.Content, storedContent.ValueString()) {
		definition.Content = storedContent.ValueString()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

// exportMatchesContent reports whether every property of the exported definition is present in the
// stored content with the same value. Properties only the stored content has are ignored, since the
// backend no longer exports some properties older exports carry (such as state, estimated_duration,
// variant_sandboxes or a level's order) and ignores them on import. Arrays must match element by
// element. Invalid JSON on either side never matches.
func exportMatchesContent(exported, stored string) bool {
	var exportedValue, storedValue any

	if err := json.Unmarshal([]byte(exported), &exportedValue); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(stored), &storedValue); err != nil {
		return false
	}

	return exportedSubsetOf(exportedValue, storedValue)
}

func exportedSubsetOf(exported, stored any) bool {
	switch exportedValue := exported.(type) {
	case map[string]any:
		storedValue, ok := stored.(map[string]any)
		if !ok {
			return false
		}
		for key, exportedProperty := range exportedValue {
			storedProperty, ok := storedValue[key]
			if !ok || !exportedSubsetOf(exportedProperty, storedProperty) {
				return false
			}
		}
		return true
	case []any:
		storedValue, ok := stored.([]any)
		if !ok || len(storedValue) != len(exportedValue) {
			return false
		}
		for i := range exportedValue {
			if !exportedSubsetOf(exportedValue[i], storedValue[i]) {
				return false
			}
		}
		return true
	default:
		return exported == stored
	}
}

func (r *trainingDefinitionResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *trainingDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.Int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTrainingDefinition(ctx, id.ValueInt64())
	if errors.Is(err, crczp.ErrNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete training definition, got error: %s", err))
		return
	}
}

func (r *trainingDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import training definition, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

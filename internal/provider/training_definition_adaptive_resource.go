package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"strconv"
	"terraform-provider-crczp/internal/plan_modifiers"

	"github.com/cyberrangecz/go-client/pkg/crczp"
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
var _ resource.Resource = &trainingDefinitionAdaptiveResource{}
var _ resource.ResourceWithImportState = &trainingDefinitionAdaptiveResource{}
var _ resource.ResourceWithConfigure = &trainingDefinitionAdaptiveResource{}

func NewTrainingDefinitionAdaptiveResource() resource.Resource {
	return &trainingDefinitionAdaptiveResource{}
}

// trainingDefinitionAdaptiveResource defines the resource implementation.
type trainingDefinitionAdaptiveResource struct {
	client *crczp.Client
}

func (r *trainingDefinitionAdaptiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_training_definition_adaptive"
}

func (r *trainingDefinitionAdaptiveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Adaptive training definition",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Id of the training definition",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "JSON with the exported training definition",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					plan_modifiers.JSONNormalizePlanModifier{},
				},
			},
		},
	}
}

func (r *trainingDefinitionAdaptiveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *trainingDefinitionAdaptiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var content string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("content"), &content)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	definition, err := r.client.CreateTrainingDefinitionAdaptive(ctx, content)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create training definition adaptive, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("created training definition adaptive %d", definition.Id))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *trainingDefinitionAdaptiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id types.Int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	definition, err := r.client.GetTrainingDefinitionAdaptive(ctx, id.ValueInt64())
	if errors.Is(err, crczp.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read training definition adaptive, got error: %s", err))
		return
	}

	var diagnostics diag.Diagnostics
	definition.Content, diagnostics = plan_modifiers.NormalizeJSON(definition.Content)
	if diagnostics != nil {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *trainingDefinitionAdaptiveResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *trainingDefinitionAdaptiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.Int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	err := r.client.DeleteTrainingDefinitionAdaptive(ctx, id.ValueInt64())
	if errors.Is(err, crczp.ErrNotFound) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete training definition adaptive, got error: %s", err))
		return
	}
}

func (r *trainingDefinitionAdaptiveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to import training definition adaptive, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

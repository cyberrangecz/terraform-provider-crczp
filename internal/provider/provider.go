package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/cyberrangecz/go-client/pkg/crczp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure CrczpProvider satisfies various provider interfaces.
var _ provider.Provider = &CrczpProvider{}

// CrczpProvider defines the provider implementation.
type CrczpProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CrczpProviderModel describes the provider data model.
type CrczpProviderModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Token      types.String `tfsdk:"token"`
	ClientID   types.String `tfsdk:"client_id"`
	RetryCount types.Int64  `tfsdk:"retry_count"`
}

func (p *CrczpProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "crczp"
	resp.Version = p.version
}

func (p *CrczpProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "URI of the homepage of the CRCZP instance, like `https://my.crczp.instance.ex`. Can be set with `CRCZP_ENDPOINT` environmental variable.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "`username` of the user to login as with `password`. Use either `username` and `password` or just `token`. Can be set with `CRCZP_USERNAME` environmental variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "`password` of the user to login as with `username`. Use either `username` and `password` or just `token`. Can be set with `CRCZP_PASSWORD` environmental variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Bearer token to be used. Takes precedence before `username` and `password`. Bearer tokens usually have limited lifespan. Can be set with `CRCZP_TOKEN` environmental variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "CRCZP local OIDC client ID. Will be ignored when `token` is set. Defaults to `CRCZP-Client`. Can be set with `CRCZP_CLIENT_ID` environmental variable. See [how to get CRCZP client_id](https://registry.terraform.io/vydrazde/crczp/latest/docs/guides/getting_oidc_client_id).",
				Optional:            true,
			},
			"retry_count": schema.Int64Attribute{
				MarkdownDescription: "How many times to retry failed HTTP requests. There is a delay of 100ms before the first retry. For each following retry, the delay is doubled. Defaults to 0. Can be set with `CRCZP_RETRY_COUNT` environmental variable.",
				Optional:            true,
			},
		},
	}
}

func (p *CrczpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CrczpProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown CRCZP API Endpoint",
			"The provider cannot create the CRCZP API client as there is an unknown configuration value for the CRCZP API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRCZP_ENDPOINT environment variable.",
		)
	}
	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown CRCZP API Username",
			"The provider cannot create the CRCZP API client as there is an unknown configuration value for the CRCZP API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRCZP_USERNAME environment variable.",
		)
	}
	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown CRCZP API Password",
			"The provider cannot create the CRCZP API client as there is an unknown configuration value for the CRCZP API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRCZP_PASSWORD environment variable.",
		)
	}
	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown CRCZP API Token",
			"The provider cannot create the CRCZP API client as there is an unknown configuration value for the CRCZP API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRCZP_TOKEN environment variable.",
		)
	}
	if data.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown CRCZP API Client ID",
			"The provider cannot create the CRCZP API client as there is an unknown configuration value for the CRCZP API client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CRCZP_CLIENT_ID environment variable.",
		)
	}

	endpoint := os.Getenv("CRCZP_ENDPOINT")
	username := os.Getenv("CRCZP_USERNAME")
	password := os.Getenv("CRCZP_PASSWORD")
	token := os.Getenv("CRCZP_TOKEN")
	clientId := os.Getenv("CRCZP_CLIENT_ID")
	retryCountStr := os.Getenv("CRCZP_RETRY_COUNT")

	retryCount, err := strconv.Atoi(retryCountStr)
	if err != nil {
		retryCount = 0
	}

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}
	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}
	if !data.ClientID.IsNull() {
		clientId = data.ClientID.ValueString()
	}
	if !data.RetryCount.IsNull() && !data.RetryCount.IsUnknown() {
		retryCount = int(data.RetryCount.ValueInt64())
	}

	if clientId == "" {
		clientId = "CRCZP-Client"
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing CRCZP API Endpoint",
			"The provider cannot create the CRCZP API client as there is a missing or empty value for the CRCZP API endpoint. "+
				"Set the host value in the configuration or use the CRCZP_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if token == "" && (username == "" || password == "") {
		resp.Diagnostics.AddError(
			"Missing CRCZP API Token or Username and Password",
			"The provider cannot create the CRCZP API client as there is a missing or empty value for the CRCZP API token or username and password. "+
				"Set the host value in the configuration or use the CRCZP_TOKEN, CRCZP_USERNAME and CRCZP_PASSWORD environment variables. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "crczp_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "crczp_username", username)
	ctx = tflog.SetField(ctx, "crczp_password", password)
	ctx = tflog.SetField(ctx, "crczp_token", token)
	ctx = tflog.SetField(ctx, "client_id", clientId)
	ctx = tflog.SetField(ctx, "retry_count", retryCount)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "crczp_password", "crczp_token")

	tflog.Debug(ctx, "Creating CRCZP client")
	var client *crczp.Client

	if token != "" {
		client, err = crczp.NewClientWithToken(endpoint, clientId, token)
	} else {
		client, err = crczp.NewClient(endpoint, clientId, username, password)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CRCZP API Client",
			"An unexpected error occurred when creating the CRCZP API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"CRCZP Client Error: "+err.Error(),
		)
		return
	}
	client.RetryCount = retryCount
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured CRCZP client", map[string]any{"success": true})
}

func (p *CrczpProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSandboxDefinitionResource,
		NewSandboxPoolResource,
		NewSandboxAllocationUnitResource,
		NewTrainingDefinitionResource,
		NewTrainingDefinitionAdaptiveResource,
	}
}

func (p *CrczpProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSandboxRequestOutputDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CrczpProvider{
			version: version,
		}
	}
}

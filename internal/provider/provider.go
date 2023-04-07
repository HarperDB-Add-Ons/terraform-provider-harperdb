package provider

import (
	"context"
	"fmt"

	harperdb "github.com/HarperDB-Add-Ons/sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure HarperDBProvider satisfies various provider interfaces.
var _ provider.Provider = &HarperDBProvider{}

// HarperDBProvider defines the provider implementation.
type HarperDBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HarperDBProviderModel describes the provider data model.
type HarperDBProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *HarperDBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "harperdb"
	resp.Version = p.version
}

func (p *HarperDBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "HarperDB endpoint",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "HarperDB super-user username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "HarperDB super-user password",
				Optional:            true,
			},
		},
	}
}

func (p *HarperDBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HarperDBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := harperdb.NewClient(data.Endpoint.ValueString(), data.Username.ValueString(), data.Password.ValueString())
	tflog.Error(ctx, fmt.Sprintf("%+v", client))
	// client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *HarperDBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
		NewRoleResource,
		NewTableResource,
		NewUserResource,
	}
}

func (p *HarperDBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HarperDBProvider{
			version: version,
		}
	}
}

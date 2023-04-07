package provider

import (
	"context"
	"fmt"

	harperdb "github.com/HarperDB-Add-Ons/sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TableResource{}
var _ resource.ResourceWithImportState = &TableResource{}

func NewTableResource() resource.Resource {
	return &TableResource{}
}

// TableResource defines the resource implementation.
type TableResource struct {
	client *harperdb.Client
}

// TableResourceModel describes the resource data model.
type TableResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Schema        types.String `tfsdk:"schema"`
	Name          types.String `tfsdk:"name"`
	HashAttribute types.String `tfsdk:"hash_attribute"`
}

func (r *TableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_table"
}

func (r *TableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Table resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the Table, <schema.name>",
				Computed:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "Default (starting) Role",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Username",
				Required:            true,
			},
			"hash_attribute": schema.StringAttribute{
				MarkdownDescription: "User password",
				Required:            true,
			},
		},
	}
}

func (r *TableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*harperdb.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *harperdb.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TableResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema := data.Schema.ValueString()
	name := data.Name.ValueString()
	hashAttribute := data.HashAttribute.ValueString()
	err := r.client.CreateTable(schema, name, hashAttribute)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create User, got error: %s", err))
		return
	}

	// No internal ID is exposed for users
	data.ID = types.StringValue(fmt.Sprintf("%s.%s", schema, name))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a User resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *TableResourceModel

	// Read Terraform prior state data into the model
	// The id and the name are the only configurable elements.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("updating User unsupported", "")
	return
}

func (r *TableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TableResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema := data.Schema.ValueString()
	name := data.Name.ValueString()
	hashAttribute := data.HashAttribute.ValueString()
	err := r.client.DropTable(schema, name, hashAttribute)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to drop User, got error: %s", err))
		return
	}
}

func (r *TableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

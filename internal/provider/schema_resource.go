package provider

import (
	"context"
	"fmt"

	harperdb "github.com/HarperDB-Add-Ons/sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SchemaResource{}

// var _ resource.ResourceWithImportState = &SchemaResource{}

func NewSchemaResource() resource.Resource {
	return &SchemaResource{}
}

// SchemaResource defines the resource implementation.
type SchemaResource struct {
	client *harperdb.Client
}

// SchemaResourceModel describes the resource data model.
type SchemaResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *SchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Schema name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the schema",
				Computed:            true,
			},
		},
	}
}

func (r *SchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateSchema(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	data.ID = data.Name

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a schema resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaResourceModel

	// Read Terraform prior state data into the model
	// The id and the name are the only configurable elements.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("updating schema unsupported", "")
	return
}

func (r *SchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DropSchema(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to drop schema, got error: %s", err))
		return
	}
}

func (r *SchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

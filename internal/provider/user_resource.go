package provider

import (
	"context"
	"fmt"

	harperdb "github.com/HarperDB-Add-Ons/sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *harperdb.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Role     types.String `tfsdk:"role"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Active   types.Bool   `tfsdk:"active"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the User",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Default (starting) Role",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "User password",
				Required:            true,
				Sensitive:           true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Is the account active",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddUser(data.Username.ValueString(), data.Password.ValueString(), data.Role.ValueString(), data.Active.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create User, got error: %s", err))
		return
	}

	// No internal ID is exposed for users
	data.ID = data.Username

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a User resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	// The id and the name are the only configurable elements.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// We only need to preserve the ID.
	err := r.client.AlterUser(data.Username.ValueString(), data.Password.ValueString(), data.Role.ValueString(), data.Active.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update User, got error: %s", err))
		return
	}

	data.ID = data.Username

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DropUser(data.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to drop User, got error: %s", err))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

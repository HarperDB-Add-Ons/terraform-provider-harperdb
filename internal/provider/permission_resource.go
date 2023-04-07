package provider

import (
	"context"
	"fmt"

	harperdb "github.com/HarperDB-Add-Ons/sdk-go"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PermissionResource{}
var _ resource.ResourceWithImportState = &PermissionResource{}

func NewPermissionResource() resource.Resource {
	return &PermissionResource{}
}

// PermissionResource defines the resource implementation.
type PermissionResource struct {
	client *harperdb.Client
}

// PermissionResourceModel describes the resource data model.
type PermissionResourceModel struct {
	ID          types.String `tfsdk:"id"` // Derived from the resource-name
	SuperUser   types.Bool   `tfsdk:"super_user"`
	ClusterUser types.Bool   `tfsdk:"cluster_user"`
	// SchemaPermissions types.Map    `tfsdk:"schema_permissions"`
	TablePermissions types.Map `tfsdk:"table_permissions"`
}

func (r *PermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *PermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Permission resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the Permission",
				Computed:            true,
			},
			"super_user": schema.BoolAttribute{
				MarkdownDescription: "Is super user",
				Optional:            true,
			},
			"cluster_user": schema.BoolAttribute{
				MarkdownDescription: "is cluster user",
				Optional:            true,
			},
			// "schema_permissions": schema.MapNestedAttribute{
			// 	MarkdownDescription: "Schema Permissions",
			// 	Optional:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"read": schema.BoolAttribute{
			// 				Optional: true,
			// 				Computed: true,
			// 				Default:  booldefault.StaticBool(false),
			// 			},
			// 			"insert": schema.BoolAttribute{
			// 				Optional: true,
			// 				Computed: true,
			// 				Default:  booldefault.StaticBool(false),
			// 			},
			// 			"update": schema.BoolAttribute{
			// 				Optional: true,
			// 				Computed: true,
			// 				Default:  booldefault.StaticBool(false),
			// 			},
			// 			"delete": schema.BoolAttribute{
			// 				Optional: true,
			// 				Computed: true,
			// 				Default:  booldefault.StaticBool(false),
			// 			},
			// 		},
			// 	},
			// },
			"table_permissions": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"read": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"insert": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"update": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"delete": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
				Optional: true,
			},
			// "table_permissions": schema.MapNestedAttribute{
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"testing": schema.BoolAttribute{
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// 	// AttributeTypes: map[string]attr.Type{
			// 	// 	"testing": types.BoolType,
			// 	// },

			// 	MarkdownDescription: "Table Permissions",
			// 	Optional:            true,
			// },
		},
	}
}

func (r *PermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PermissionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// No internal ID is exposed for Permission so we generate one.
	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Permission, got error: %s", err))
	}

	data.ID = basetypes.NewStringValue(id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a Permission resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PermissionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PermissionResourceModel
	var old_data *PermissionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &old_data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// We only need to preserve the ID.
	data.ID = old_data.ID

	// We simply update as we don't need to sync changes.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *PermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PermissionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// This is a state-only resource. It doesn't have a direct analogy in HarperDB.
}

func (r *PermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

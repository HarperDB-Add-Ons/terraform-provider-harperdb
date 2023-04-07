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
var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}

func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// RoleResource defines the resource implementation.
type RoleResource struct {
	client *harperdb.Client
}

// RoleResourceModel describes the resource data model.
type RoleResourceModel struct {
	ID                types.String `tfsdk:"id"`   // Derived from the resource-name
	Name              types.String `tfsdk:"name"` // Role name
	SuperUser         types.Bool   `tfsdk:"super_user"`
	ClusterUser       types.Bool   `tfsdk:"cluster_user"`
	SchemaPermissions types.Map    `tfsdk:"schema_permissions"`
	// TablePermissions  types.Map    `tfsdk:"table_permissions"`
}

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tableSchema := schema.MapNestedAttribute{
		Optional: true,
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
				"attribute_permissions": schema.ListNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required: true,
							},
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
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Permission resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the Role",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Role",
				Required:            true,
			},
			"super_user": schema.BoolAttribute{
				MarkdownDescription: "Is super user",
				Optional:            true,
			},
			"cluster_user": schema.BoolAttribute{
				MarkdownDescription: "is cluster user",
				Optional:            true,
			},
			"schema_permissions": schema.MapNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tables": tableSchema,
					},
				},
			},
		},
	}
}

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RoleResource) constructPermission(data *RoleResourceModel) harperdb.Permission {
	perm := harperdb.Permission{}
	perm.SetClusterUser(data.ClusterUser.ValueBool())
	perm.SetSuperUser(data.SuperUser.ValueBool())

	schemas := data.SchemaPermissions.Elements()

	for name, schema := range schemas {
		tables := map[string]harperdb.TablePermission{}

		tablesRaw := schema.(types.Object).Attributes()["tables"].(types.Map).Elements()
		for tname, traw := range tablesRaw {
			attr := traw.(types.Object).Attributes()
			var attributePermissions []harperdb.AttributePermissions
			attributes := attr["attribute_permissions"].(types.List).Elements()
			for _, attribute := range attributes {
				a := attribute.(types.Object).Attributes()
				attributePermissions = append(attributePermissions, harperdb.AttributePermissions{
					AttributeName: a["name"].(types.String).ValueString(),
					Read:          a["read"].(types.Bool).ValueBool(),
					Insert:        a["insert"].(types.Bool).ValueBool(),
					Update:        a["update"].(types.Bool).ValueBool(),
				})

			}
			tables[tname] = harperdb.TablePermission{
				Read:                 attr["read"].(types.Bool).ValueBool(),
				Insert:               attr["insert"].(types.Bool).ValueBool(),
				Update:               attr["update"].(types.Bool).ValueBool(),
				Delete:               attr["delete"].(types.Bool).ValueBool(),
				AttributePermissions: attributePermissions,
			}
		}

		schemaPermission := harperdb.SchemaPermission{
			Tables: tables,
		}
		perm.AddSchemaPermission(name, schemaPermission)
	}

	return perm
}

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleName := data.Name.ValueString()
	perm := r.constructPermission(data)

	role, err := r.client.AddRole(roleName, perm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create role, got error: %s !! %+v", err, role))
		return
	}
	data.ID = types.StringValue(role.ID)
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a Role resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *RoleResourceModel
	var old_data *RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &old_data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	perm := r.constructPermission(data)

	role, err := r.client.AlterRole(old_data.ID.ValueString(), data.Name.ValueString(), perm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create role, got error: %s !! %+v", err, role))
		return
	}

	// We only need to preserve the ID.
	data.ID = types.StringValue(role.ID)

	// We simply update as we don't need to sync changes.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	err := r.client.DropRole(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to drop User, got error: %s", err))
		return
	}

	// This is a state-only resource. It doesn't have a direct analogy in HarperDB.
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

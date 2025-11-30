package orgserviceaccount

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const resourceName = "org_service_account"

var _ resource.ResourceWithConfigure = &rs{}
var _ resource.ResourceWithImportState = &rs{}

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{
			ResourceName: resourceName,
		},
	}
}

type rs struct {
	config.RSCommon
}

func (r *rs) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	atlasReq, diag := NewAtlasReq(ctx, &plan)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	orgID := plan.OrgId.ValueString()

	connV2 := r.Client.AtlasV2
	orgServiceAccountReq, _, err := connV2.ServiceAccountsApi.CreateOrgServiceAccount(ctx, orgID, atlasReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	newOrgServiceAccountModel, diags := NewTFModel(ctx, orgServiceAccountReq)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newOrgServiceAccountModel)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrgId.ValueString()
	saID := state.ClientId.ValueString()

	connV2 := r.Client.AtlasV2
	getOrgServiceAccountReq, apiResp, err := connV2.ServiceAccountsApi.GetOrgServiceAccount(ctx, orgID, saID).Execute()

	if err != nil {
		if validate.StatusNotFound(apiResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	newOrgServiceAccountModel, diags := NewTFModel(ctx, getOrgServiceAccountReq)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newOrgServiceAccountModel)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	atlasUpdateReq, diag := NewAtlasUpdateReq(ctx, &plan)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	orgID := plan.OrgId.ValueString()
	saID := plan.ClientId.ValueString()

	connV2 := r.Client.AtlasV2

	orgServiceAccountReq, _, err := connV2.ServiceAccountsApi.UpdateOrgServiceAccount(ctx, saID, orgID, atlasUpdateReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	newOrgServiceAccountModel, diags := NewTFModel(ctx, orgServiceAccountReq)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newOrgServiceAccountModel)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID := state.OrgId.ValueString()
	saID := state.ClientId.ValueString()

	connV2 := r.Client.AtlasV2
	if _, err := connV2.ServiceAccountsApi.DeleteOrgServiceAccount(ctx, saID, orgID).Execute(); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	orgID, saID, err := splitImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_id"), orgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client_id"), saID)...)
}

func splitImportID(id string) (orgID, resourcePolicyID string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("use the format {org_id}-{org_service_account_id}")
		return
	}

	orgID = parts[1]
	resourcePolicyID = parts[2]
	return
}

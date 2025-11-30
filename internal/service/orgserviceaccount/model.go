package orgserviceaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

func NewTFModel(ctx context.Context, apiResp *admin.OrgServiceAccount) (*TFModel, diag.Diagnostics) {
	roles, diags := types.SetValueFrom(ctx, types.StringType, *(apiResp.Roles))
	if diags.HasError() {
		return nil, diags
	}

	secrets := NewTFSecrets(ctx, apiResp.Secrets)

	return &TFModel{
		ClientId:    types.StringPointerValue(apiResp.ClientId),
		CreatedAt:   types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedAt)),
		Description: types.StringPointerValue(apiResp.Description),
		Name:        types.StringPointerValue(apiResp.Name),
		Roles:       roles,
		Secrets:     secrets,
	}, nil
}

func NewTFSecrets(ctx context.Context, input *[]admin.ServiceAccountSecret) []TFSecretsModel {
	var nilPointer *[]admin.ServiceAccountSecret
	if input == nilPointer {
		return nil
	}
	tfSecrets := make([]TFSecretsModel, len(*input))
	for i, item := range *input {
		tfSecrets[i] = TFSecretsModel{
			CreatedAt:         types.StringPointerValue(conversion.TimePtrToStringPtr(&item.CreatedAt)),
			ExpiresAt:         types.StringPointerValue(conversion.TimePtrToStringPtr(&item.ExpiresAt)),
			Id:                types.StringPointerValue(&item.Id),
			LastUsedAt:        types.StringPointerValue(conversion.TimePtrToStringPtr(item.LastUsedAt)),
			MaskedSecretValue: types.StringPointerValue(item.MaskedSecretValue),
			Secret:            types.StringPointerValue(item.Secret),
		}
	}
	return tfSecrets
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.OrgServiceAccountRequest, diag.Diagnostics) {
	var roles []string
	diags := plan.Roles.ElementsAs(ctx, &roles, false)
	if diags.HasError() {
		return nil, diags
	}

	return &admin.OrgServiceAccountRequest{
		Name:                    plan.Name.ValueString(),
		Description:             plan.Description.ValueString(),
		SecretExpiresAfterHours: int(plan.SecretExpiresAfterHours.ValueInt64()),
		Roles:                   roles,
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.OrgServiceAccountUpdateRequest, diag.Diagnostics) {
	var roles []string
	diags := plan.Roles.ElementsAs(ctx, &roles, false)
	if diags.HasError() {
		return nil, diags
	}

	return &admin.OrgServiceAccountUpdateRequest{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		Roles:       &roles,
	}, nil
}

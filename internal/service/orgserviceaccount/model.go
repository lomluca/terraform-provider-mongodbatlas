package orgserviceaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

// TODO: `ctx` parameter and `diags` return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
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
		return []TFSecretsModel{}
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

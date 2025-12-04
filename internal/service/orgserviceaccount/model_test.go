package orgserviceaccount_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/orgserviceaccount"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

const (
	clientID                = "mdb_sa_id_a9b3c7d2e6f1g4h8i0j5k2l3"
	createdAt               = "2025-11-27T13:40:47Z"
	description             = "Service account to test the tf resource"
	name                    = "sa-test-org-owner"
	secretExpiresAt         = "2026-02-25T13:40:47Z"     // #nosec G101
	secretID                = "a9b3c7d2e6f1g4h8i0j5k2l2" // #nosec G101
	secretLastUsedAt        = "2025-11-27T08:51:17Z"     // #nosec G101
	secretMaskedSecretValue = "mdb_sa_sk_...L12x"        // #nosec G101
)

var (
	//go:embed testdata/full.json
	fullJSON string
	//go:embed testdata/no_secrets.json
	noSecretsJSON string
	//go:embed testdata/new_request.json
	newRequestJSON string
	//go:embed testdata/new_update_request.json
	newUpdateRequestJSON string
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.OrgServiceAccount
	expectedTFModel *orgserviceaccount.TFModel
}

func parseSDKModel(t *testing.T, sdkRespJSON string) *admin.OrgServiceAccount {
	t.Helper()
	var SDKModel admin.OrgServiceAccount
	err := json.Unmarshal([]byte(sdkRespJSON), &SDKModel)
	if err != nil {
		t.Fatalf("failed to unmarshal sdk response: %s", err)
	}
	return &SDKModel
}

func parseSDKRequest(t *testing.T, sdkReqJSON string) *admin.OrgServiceAccountRequest {
	t.Helper()
	var SDKRequest admin.OrgServiceAccountRequest
	err := json.Unmarshal([]byte(sdkReqJSON), &SDKRequest)
	if err != nil {
		t.Fatalf("failed to unmarshal sdk response: %s", err)
	}
	return &SDKRequest
}

func parseSDKUpdateRequest(t *testing.T, sdkReqJSON string) *admin.OrgServiceAccountUpdateRequest {
	t.Helper()
	var SDKRequest admin.OrgServiceAccountUpdateRequest
	err := json.Unmarshal([]byte(sdkReqJSON), &SDKRequest)
	if err != nil {
		t.Fatalf("failed to unmarshal sdk response: %s", err)
	}
	return &SDKRequest
}

func buildTFModel(ctx context.Context, addSecrets bool) *orgserviceaccount.TFModel {
	roles, _ := types.SetValueFrom(ctx, types.StringType, []string{"GROUP_OWNER"})
	tfModel := orgserviceaccount.TFModel{
		ClientId:    types.StringValue(clientID),
		CreatedAt:   types.StringValue(createdAt),
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		Roles:       roles,
	}

	if addSecrets == true {
		tfModel.Secrets, _ = types.SetValueFrom(ctx, orgserviceaccount.SecretObjectType, []orgserviceaccount.TFSecretsModel{
			{
				CreatedAt:         types.StringValue(createdAt),
				ExpiresAt:         types.StringValue(secretExpiresAt),
				Id:                types.StringValue(secretID),
				LastUsedAt:        types.StringValue(secretLastUsedAt),
				MaskedSecretValue: types.StringValue(secretMaskedSecretValue),
			},
		})
	}

	return &tfModel
}

func TestOrgServiceAccountSDKToTFModel(t *testing.T) {
	ctx := t.Context()

	simpleFullResponse := parseSDKModel(t, fullJSON)
	simpleFullModel := buildTFModel(ctx, true)

	noSecretsResponse := parseSDKModel(t, noSecretsJSON)
	noSecretsModel := buildTFModel(ctx, false)

	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp:         simpleFullResponse,
			expectedTFModel: simpleFullModel,
		},
		"No secrets": {
			SDKResp:         noSecretsResponse,
			expectedTFModel: noSecretsModel,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := orgserviceaccount.NewTFModel(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func TestNewAtlasReq(t *testing.T) {
	newRequestSDK := parseSDKRequest(t, newRequestJSON)
	roles, _ := types.SetValueFrom(t.Context(), types.StringType, []string{"ORG_MEMBER"})

	testCases := map[string]struct {
		tfModel        *orgserviceaccount.TFModel
		expectedSDKReq *admin.OrgServiceAccountRequest
	}{
		"Create": {
			tfModel: &orgserviceaccount.TFModel{
				Name:                    types.StringValue(name),
				SecretExpiresAfterHours: types.Int64Value(8),
				Description:             types.StringValue(description),
				Roles:                   roles,
			},
			expectedSDKReq: newRequestSDK,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := orgserviceaccount.NewAtlasReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

func TestNewAtlasUpdateReq(t *testing.T) {
	newUpdateRequestSDK := parseSDKUpdateRequest(t, newUpdateRequestJSON)
	roles, _ := types.SetValueFrom(t.Context(), types.StringType, []string{"ORG_MEMBER"})

	testCases := map[string]struct {
		tfModel        *orgserviceaccount.TFModel
		expectedSDKReq *admin.OrgServiceAccountUpdateRequest
	}{
		"Update": {
			tfModel: &orgserviceaccount.TFModel{
				Name:        types.StringValue(name),
				Description: types.StringValue(description),
				Roles:       roles,
			},
			expectedSDKReq: newUpdateRequestSDK,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := orgserviceaccount.NewAtlasUpdateReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

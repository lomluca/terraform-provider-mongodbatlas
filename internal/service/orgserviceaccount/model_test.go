package orgserviceaccount_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/orgserviceaccount"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.OrgServiceAccount
	expectedTFModel *orgserviceaccount.TFModel
}

func TestOrgServiceAccountSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp:         &admin.OrgServiceAccount{},
			expectedTFModel: &orgserviceaccount.TFModel{},
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

/*
type tfToSDKModelTestCase struct {
	tfModel        *orgserviceaccount.TFModel
	expectedSDKReq *admin.OrgServiceAccount
}


func TestOrgServiceAccountTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel:        &orgserviceaccount.TFModel{},
			expectedSDKReq: &admin.OrgServiceAccount{},
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
}*/

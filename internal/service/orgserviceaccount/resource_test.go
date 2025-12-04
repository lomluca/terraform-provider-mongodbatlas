package orgserviceaccount_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_org_service_account.test"
const dataSourceName = "data.mongodbatlas_org_service_account.test"

// TODO: if acceptance test will be run in an existing CI group of resources, the name should include the group in the prefix followed by the name of the resource e.i. TestAccStreamRSStreamInstance_basic
// In addition, if acceptance test contains testing of both resource and data sources, the RS/DS can be omitted.
func TestAccOrgServiceAccount_basic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	name := acc.RandomName()
	roles := []string{"GROUP_OWNER"}
	updatedRoles := []string{"GROUP_OWNER", "ORG_MEMBER"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyOrgServiceAccount,
		Steps: []resource.TestStep{ // TODO: verify updates and import in case of resources
			{
				Config: orgServiceAccountConfig(orgID, name, description, roles),
				Check:  orgServiceAccountAttributeChecks(name, description, roles),
			},
			{
				Config: orgServiceAccountConfig(orgID, name, description, updatedRoles),
				Check:  orgServiceAccountAttributeChecks(name, description, updatedRoles),
			},
			{
				ResourceName:                         resourceName,
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportStateIdFunc:                    checkOrgServiceAccountImportStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func orgServiceAccountConfig(orgID, name, description string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		resource "mongodbatlas_org_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = %[3]q
			roles                      = [%[4]q]
			secret_expires_after_hours = 12
		}
			
		data "mongodbatlas_org_service_account" "test" {
			org_id    = %[1]q
			client_id = mongodbatlas_org_service_account.test.client_id
		}
	`, orgID, name, description, rolesStr)
}

func orgServiceAccountAttributeChecks(name, description string, roles []string) resource.TestCheckFunc {
	attrsSet := []string{"team_id"}
	attrsMap := map[string]string{
		"name":        name,
		"description": description,
		"roles.#":     fmt.Sprint(len(roles)),
	}
	extraChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
	}
	for _, role := range roles {
		extraChecks = append(extraChecks, resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", role))
	}

	return acc.CheckRSAndDS(resourceName, nil, nil, attrsSet, attrsMap)
}

func checkDestroyOrgServiceAccount(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_org_service_account" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		orgServiceAccountID := rs.Primary.Attributes["client_id"]
		conn := acc.ConnV2()
		orgServiceAccountListResp, _, err := conn.ServiceAccountsApi.ListOrgServiceAccounts(
			context.Background(),
			orgID,
		).Execute()
		if err != nil {
			continue
		}

		if orgServiceAccountListResp != nil && orgServiceAccountListResp.Results != nil {
			results := *orgServiceAccountListResp.Results
			for i := range results {
				if *results[i].ClientId == orgServiceAccountID {
					return fmt.Errorf("Org Service Account %s still exists", orgServiceAccountID)
				}
			}
		}
	}
	return nil
}

func checkOrgServiceAccountImportStateIDFunc(resourceName string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrs := s.RootModule().Resources[resourceName].Primary.Attributes
		return attrs["client_id"], nil
	}
}

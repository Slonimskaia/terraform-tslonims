package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		configResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(configResponseJSON), &configResponse)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createRequestJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")
		createResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyCreateResponse.json"))
		createResponse := appsec.CreateCustomDenyResponse{}
		json.Unmarshal([]byte(createResponseJSON), &createResponse)
		client.On("CreateCustomDeny",
			mock.Anything,
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: createRequestJSON},
		).Return(&createResponse, nil)

		getResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyGetResponse.json"))
		getResponse := appsec.GetCustomDenyResponse{}
		json.Unmarshal([]byte(getResponseJSON), &getResponse)
		client.On("GetCustomDeny",
			mock.Anything,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil).Times(3)

		updateRequestJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheFalse.json")
		updateResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyUpdateResponse.json"))
		updateResponse := appsec.UpdateCustomDenyResponse{}
		json.Unmarshal([]byte(updateResponseJSON), &updateResponse)
		client.On("UpdateCustomDeny",
			mock.Anything,
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918", JsonPayloadRaw: updateRequestJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdateJSON := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyGetResponseAfterUpdate.json"))
		getResponseAfterUpdate := appsec.GetCustomDenyResponse{}
		json.Unmarshal([]byte(getResponseAfterUpdateJSON), &getResponseAfterUpdate)
		client.On("GetCustomDeny",
			mock.Anything,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		removeResponse := appsec.RemoveCustomDenyResponse{}
		json.Unmarshal([]byte(removeResponseJSON), &removeResponse)
		client.On("RemoveCustomDeny",
			mock.Anything,
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
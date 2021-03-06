package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiWAFProtections_data_basic(t *testing.T) {
	t.Run("match by WAFProtections ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetWAFProtectionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSWAFProtections/WAFProtections.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetWAFProtections",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAFProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyApplicationLayerControls: false},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAFProtections/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_waf_protection.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

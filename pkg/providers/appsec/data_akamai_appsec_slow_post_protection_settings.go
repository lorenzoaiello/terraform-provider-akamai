package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSlowPostProtectionSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlowPostProtectionSettingsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSlowPostProtectionSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingsRead")
	//CorrelationID := "[APPSEC][resourceSlowPostProtectionSettings-" + meta.OperationID() + "]"

	getSlowPostProtectionSettings := v2.GetSlowPostProtectionSettingsRequest{}

	getSlowPostProtectionSettings.ConfigID = d.Get("config_id").(int)
	getSlowPostProtectionSettings.Version = d.Get("version").(int)
	getSlowPostProtectionSettings.PolicyID = d.Get("policy_id").(string)

	slowpostprotectionsettings, err := client.GetSlowPostProtectionSettings(ctx, getSlowPostProtectionSettings)
	if err != nil {
		logger.Warnf("calling 'getSlowPostProtectionSettings': %s", err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "slowPostDS", slowpostprotectionsettings)
	//edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("slowPost outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getSlowPostProtectionSettings.ConfigID))

	return nil
}
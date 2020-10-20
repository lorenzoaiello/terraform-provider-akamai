package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRuleActionsRead,
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

func dataSourceCustomRuleActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionsRead")
	CorrelationID := "[APPSEC][resourceCustomRuleActions-" + meta.OperationID() + "]"

	getCustomRuleActions := v2.GetCustomRuleActionsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleActions.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleActions.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleActions.PolicyID = policyid

	customruleactions, err := client.GetCustomRuleActions(ctx, getCustomRuleActions)
	if err != nil {
		logger.Warnf("calling 'getCustomRuleActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customRuleAction", customruleactions)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("customRuleAction outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	//d.Set("rule_id", ruleid)
	d.SetId(strconv.Itoa(getCustomRuleActions.ConfigID))

	return nil
}
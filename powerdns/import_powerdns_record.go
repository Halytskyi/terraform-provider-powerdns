package powerdns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePowerDNSRecordImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	client := meta.(*Client)

	values := strings.Split(d.Id(), idSeparator)

	if len(values) != 2 && len(values) != 3 {
		return nil, fmt.Errorf("invalid id provided, expected format: {fqdn.}%s{type}[%s{set_ptr}]", idSeparator, idSeparator)
	}

	recordID := values[0] + idSeparator + values[1]
	recordFQDN := values[0]
	recordZone := strings.TrimSuffix(strings.SplitN(recordFQDN, ".", 2)[1], ".")

	var recordSetPTR bool
	if len(values) == 3 {
		if values[2] == "true" {
			recordSetPTR = true
		} else if values[2] == "false" {
			recordSetPTR = false
		} else {
			return nil, fmt.Errorf("invalid 'set_ptr' parameter, should be 'true' or 'false'")
		}
	}

	records, err := client.ListRecordsByID(recordZone, recordID)

	if err != nil {
		return nil, fmt.Errorf("Couldn't fetch PowerDNS Record: %s", err)
	} else if len(records) == 0 {
		return nil, fmt.Errorf("Record doesn't exists in PowerDNS")
	}

	var recs []string
	recs = append(recs, records[0].Content)

	d.SetId(recordID)
	d.Set("name", records[0].Name)
	d.Set("zone", recordZone)
	d.Set("records", recs)
	d.Set("type", records[0].Type)
	d.Set("ttl", records[0].TTL)
	if recordSetPTR == true {
		d.Set("set_ptr", recordSetPTR)
	}
	results[0] = d

	return results, nil
}

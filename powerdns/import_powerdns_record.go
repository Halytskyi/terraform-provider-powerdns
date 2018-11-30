package powerdns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePowerDNSRecordImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	client := meta.(*Client)

	recordID := d.Id()
	recordFQDN := strings.Split(d.Id(), idSeparator)[0]
	recordZone := strings.TrimSuffix(strings.SplitN(recordFQDN, ".", 2)[1], ".")

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
	results[0] = d

	return results, nil
}

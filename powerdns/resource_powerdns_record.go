package powerdns

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourcePDNSRecordCreate,
		Read:   resourcePDNSRecordRead,
		Delete: resourcePDNSRecordDelete,
		Exists: resourcePDNSRecordExists,
		Importer: &schema.ResourceImporter{
			State: resourcePowerDNSRecordImportState,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
			},

			"set_ptr": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	rrSet := ResourceRecordSet{
		Name: d.Get("name").(string),
		Type: d.Get("type").(string),
		TTL:  d.Get("ttl").(int),
	}

	zone := d.Get("zone").(string)
	ttl := d.Get("ttl").(int)
	recs := d.Get("records").(*schema.Set).List()
	set_ptr := d.Get("set_ptr").(bool)

	if len(recs) > 0 {
		records := make([]Record, 0, len(recs))
		for _, recContent := range recs {
			records = append(records, Record{Name: rrSet.Name, Type: rrSet.Type, TTL: ttl,
				SetPTR: set_ptr, Content: recContent.(string)})
		}
		rrSet.Records = records

		log.Printf("[DEBUG] Creating PowerDNS Record: %#v", rrSet)

		recId, err := client.ReplaceRecordSet(zone, rrSet)
		if err != nil {
			return fmt.Errorf("Failed to create PowerDNS Record: %s", err)
		}

		d.SetId(recId)
		log.Printf("[INFO] Created PowerDNS Record with ID: %s", d.Id())

	} else {
		log.Printf("[DEBUG] Deleting empty PowerDNS Record: %#v", rrSet)
		err := client.DeleteRecordSet(zone, rrSet.Name, rrSet.Type)
		if err != nil {
			return fmt.Errorf("Failed to delete PowerDNS Record: %s", err)
		}

		d.SetId(rrSet.Id())
	}

	return resourcePDNSRecordRead(d, meta)
}

func resourcePDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	log.Printf("[DEBUG] Reading PowerDNS Record: %s", d.Id())
	records, err := client.ListRecordsByID(d.Get("zone").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't fetch PowerDNS Record: %s", err)
	}

	recs := make([]string, 0, len(records))
	for _, r := range records {
		recs = append(recs, r.Content)
	}
	d.Set("records", recs)

	if len(records) > 0 {
		d.Set("ttl", records[0].TTL)
	}

	return nil
}

func resourcePDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	is_set_ptr := d.Get("set_ptr").(bool)
	recs := d.Get("records").(*schema.Set).List()

	log.Printf("[INFO] Deleting PowerDNS Record: %s", d.Id())
	err := client.DeleteRecordSetByID(d.Get("zone").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting PowerDNS Record: %s", err)
	}

	if is_set_ptr {
		for _, ip := range recs {
			ip_octets := strings.Split(ip.(string), ".")
			ptr_name := fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", ip_octets[3], ip_octets[2], ip_octets[1], ip_octets[0])
			log.Printf("[INFO] Deleting PTR PowerDNS Record: %s", ptr_name)
			for i := 0; i < 3; i++ {
				var ptr_zone_prefix string
				if i == 0 {
					ptr_zone_prefix = fmt.Sprintf("%s.%s.%s", ip_octets[2], ip_octets[1], ip_octets[0])
				} else if i == 1 {
					ptr_zone_prefix = fmt.Sprintf("%s.%s", ip_octets[1], ip_octets[0])
				} else if i == 2 {
					ptr_zone_prefix = ip_octets[0]
				}
				ptr_zone := fmt.Sprintf("%s.in-addr.arpa", ptr_zone_prefix)
				is_ptrRecord, err := client.RecordExists(ptr_zone, ptr_name, "PTR")
				if err != nil {
					return fmt.Errorf("Error during check PTR record: %s, reason: %s", ptr_name, err)
				}
				if is_ptrRecord {
					err := client.DeleteRecordSet(ptr_zone, ptr_name, "SetPTR")
					if err != nil {
						return fmt.Errorf("Error deleting PTR record: %s, reason: %s", ptr_name, err)
					}
					break
				}
			}
		}
	}

	return nil
}

func resourcePDNSRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	tpe := d.Get("type").(string)

	log.Printf("[INFO] Checking existence of PowerDNS Record: %s, %s", name, tpe)

	client := meta.(*Client)
	exists, err := client.RecordExists(zone, name, tpe)

	if err != nil {
		return false, fmt.Errorf("Error checking PowerDNS Record: %s", err)
	} else {
		return exists, nil
	}
}

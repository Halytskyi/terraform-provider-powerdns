package powerdns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccImportPDNSRecord_A(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}

		expectedName := "redis.sysa.xyz."
		expectedZone := "sysa.xyz"
		expectedValue1 := "1.1.1.1"
		expectedValue2 := "2.2.2.2"
		expectedType := "A"
		expectedTTL := "60"
		return compareState(s[0], expectedName, expectedZone, expectedValue1, expectedValue2, expectedType, expectedTTL)
	}

	resourceName := "powerdns_record.test-a"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testPDNSRecordConfigA,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "redis.sysa.xyz.:::A",
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

func compareState(recordState *terraform.InstanceState, expectedName, expectedZone, expectedValue1, expectedValue2, expectedType, expectedTTL string) error {
	if recordState.Attributes["zone"] != expectedZone {
		return fmt.Errorf("expected zone of %s, %s received",
			expectedZone, recordState.Attributes["zone"])
	}
	if recordState.Attributes["name"] != expectedName {
		return fmt.Errorf("expected name of %s, %s received",
			expectedName, recordState.Attributes["name"])
	}
	if recordState.Attributes["records.218290772"] != expectedValue1 {
		return fmt.Errorf("expected value of %s, %s received",
			expectedValue1, recordState.Attributes["records.218290772"])
	}
	if recordState.Attributes["records.3758446074"] != expectedValue2 {
		return fmt.Errorf("expected value of %s, %s received",
			expectedValue2, recordState.Attributes["records.3758446074"])
	}
	if recordState.Attributes["type"] != expectedType {
		return fmt.Errorf("expected type of %s, %s received",
			expectedType, recordState.Attributes["type"])
	}
	if recordState.Attributes["ttl"] != expectedTTL {
		return fmt.Errorf("expected TTL of %s, %s received",
			expectedTTL, recordState.Attributes["ttl"])
	}

	return nil
}

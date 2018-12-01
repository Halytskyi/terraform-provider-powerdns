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
		expectedRecord1 := "1.1.1.1"
		expectedRecord2 := "2.2.2.2"
		expectedType := "A"
		expectedTTL := "60"
		expectedSetPTR := ""
		return compareState(s[0], expectedName, expectedZone, expectedRecord1, expectedRecord2, expectedType, expectedTTL, expectedSetPTR)
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

func TestAccImportPDNSRecord_A_WithPTR(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}

		expectedName := "redis-ptr.sysa.xyz."
		expectedZone := "sysa.xyz"
		expectedRecord1 := "1.1.1.1"
		expectedRecord2 := "2.2.2.2"
		expectedType := "A"
		expectedTTL := "60"
		expectedSetPTR := "true"
		return compareState(s[0], expectedName, expectedZone, expectedRecord1, expectedRecord2, expectedType, expectedTTL, expectedSetPTR)
	}

	resourceName := "powerdns_record.test-a-ptr"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testPDNSRecordConfigAWithPTR,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "redis-ptr.sysa.xyz.:::A:::true",
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

func compareState(recordState *terraform.InstanceState, expectedName, expectedZone, expectedRecord1, expectedRecord2, expectedType, expectedTTL string, expectedSetPTR string) error {
	if recordState.Attributes["zone"] != expectedZone {
		return fmt.Errorf("expected zone of %s, %s received",
			expectedZone, recordState.Attributes["zone"])
	}
	if recordState.Attributes["name"] != expectedName {
		return fmt.Errorf("expected name of %s, %s received",
			expectedName, recordState.Attributes["name"])
	}
	if recordState.Attributes["records.218290772"] != expectedRecord1 {
		return fmt.Errorf("expected record of %s, %s received",
			expectedRecord1, recordState.Attributes["records.218290772"])
	}
	if recordState.Attributes["records.3758446074"] != expectedRecord2 {
		return fmt.Errorf("expected record of %s, %s received",
			expectedRecord2, recordState.Attributes["records.3758446074"])
	}
	if recordState.Attributes["type"] != expectedType {
		return fmt.Errorf("expected type of %s, %s received",
			expectedType, recordState.Attributes["type"])
	}
	if recordState.Attributes["ttl"] != expectedTTL {
		return fmt.Errorf("expected TTL of %s, %s received",
			expectedTTL, recordState.Attributes["ttl"])
	}
	if recordState.Attributes["set_ptr"] != expectedSetPTR {
		return fmt.Errorf("expected set_ptr of %s, %s received",
			expectedSetPTR, recordState.Attributes["set_ptr"])
	}
	return nil
}

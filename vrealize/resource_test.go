package vrealize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/terraform-provider-vra7/utils"
)

func TestChangeValueFunction(t *testing.T) {
	request_template_original := CatalogItemRequestTemplate{}
	request_template_backup := CatalogItemRequestTemplate{}
	strJson := `{"type":"com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest","catalogItemId":"e5dd4fba-45ed-4943-b1fc-7f96239286be","requestedFor":"jason@corp.local","businessGroupId":"53619006-56bb-4788-9723-9eab79752cc1","description":null,"reasons":null,"data":{"CentOS_6.3":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*CentOS_6.3","data":{"_allocation":{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation","typeFilter":null,"data":{"machines":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Allocation.Machine","typeFilter":null,"data":{"machine_id":"","nics":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.Nic","typeFilter":null,"data":{"address":"","assignment_type":"Static","external_address":"","id":null,"load_balancing":null,"network":null,"network_profile":null}}]}}]}},"_cluster":1,"_hasChildren":false,"cpu":1,"datacenter_location":null,"description":"Basic IaaS CentOS Machine","disks":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"data":{"capacity":3,"custom_properties":null,"id":1450725224066,"initial_location":"","is_clone":true,"label":"Hard disk 1","storage_reservation_policy":"","userCreated":false,"volumeId":0}}],"guest_customization_specification":"CentOS","max_network_adapters":-1,"max_per_user":0,"max_volumes":60,"memory":512,"nics":null,"os_arch":"x86_64","os_distribution":null,"os_type":"Linux","os_version":null,"property_groups":null,"reservation_policy":null,"security_groups":[],"security_tags":[],"storage":3}},"_archiveDays":5,"_leaseDays":null,"_number_of_instances":1,"corp192168110024":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"CentOS63*corp192168110024","data":{"_hasChildren":false}}}}`
	json.Unmarshal([]byte(strJson), &request_template_original)
	json.Unmarshal([]byte(strJson), &request_template_backup)
	var flag bool

	request_template_original.Data, flag = replaceValueInRequestTemplate(request_template_original.Data, "false_field", 1000)
	if flag != false {
		t.Errorf("False value updated")
	}

	eq := reflect.DeepEqual(request_template_backup.Data, request_template_original.Data)
	if !eq {
		t.Errorf("False value updated")
	}

	request_template_original.Data, flag = replaceValueInRequestTemplate(request_template_original.Data, "storage", 1000)
	if flag == false {
		t.Errorf("Failed to update interface value")
	}

	eq2 := reflect.DeepEqual(request_template_backup.Data, request_template_original.Data)
	if eq2 {
		t.Errorf("Failed to update interface value")
	}

}

func TestConfigValidityFunction(t *testing.T) {

	mockRequestTemplate := GetMockRequestTemplate()

	// a resource_configuration map is created with valid components
	// all combinations of components name and properties are created with dots
	mockConfigResourceMap := make(map[string]interface{})
	mockConfigResourceMap["mock.test.machine1.cpu"] = 2
	mockConfigResourceMap["mock.test.machine1.mock.storage"] = 8

	resourceSchema := resourceSchema()

	resourceDataMap := map[string]interface{}{
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	readProviderConfiguration(mockResourceData)
	err := checkResourceConfigValidity(mockRequestTemplate)
	if err != nil {
		t.Errorf("The terraform config is valid, failed to validate. Expecting no error, but found %v ", err.Error())
	}

	mockConfigResourceMap["machine2.mock.cpu"] = 2
	mockConfigResourceMap["machine2.storage"] = 2

	resourceDataMap = map[string]interface{}{
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	readProviderConfiguration(mockResourceData)

	err = checkResourceConfigValidity(mockRequestTemplate)
	if err != nil {
		t.Errorf("The terraform config is valid, failed to validate. Expecting no error, but found %v ", err.Error())
	}

	mockConfigResourceMap["mock.machine3.vSphere.mock.cpu"] = 2
	resourceDataMap = map[string]interface{}{
		utils.CATALOG_ID:             "abcdefghijklmn",
		utils.RESOURCE_CONFIGURATION: mockConfigResourceMap,
	}

	mockResourceData = schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)
	readProviderConfiguration(mockResourceData)

	var mockInvalidKeys []string
	mockInvalidKeys = append(mockInvalidKeys, "mock.machine3.vSphere.mock.cpu")

	validityErr := fmt.Sprintf(utils.CONFIG_INVALID_ERROR, strings.Join(mockInvalidKeys, ", "))
	err = checkResourceConfigValidity(mockRequestTemplate)
	// this should throw an error as none of the string combinations (mock, mock.machine3, mock.machine3.vsphere, etc)
	// matches the component names(mock.test.machine1 and machine2) in the request template
	if err == nil {
		t.Errorf("The terraform config is invalid. failed to validate. Expected the error %v. but found no error", validityErr)
	}

	if err.Error() != validityErr {
		t.Errorf("Expected: %v, but Found: %v", validityErr, err.Error())
	}
}

// creates a mock request template from a request template template json file
func GetMockRequestTemplate() *CatalogItemRequestTemplate {

	ps := utils.GetPathSeparator()
	filePath := os.Getenv("GOPATH") + ps + "src" + ps + "github.com" + ps +
		"vmware" + ps + "terraform-provider-vra7" + ps + "resources" + ps + "MockRequestTemplate"

	absPath, _ := filepath.Abs(filePath)

	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	mockRequestTemplate := CatalogItemRequestTemplate{}
	json.Unmarshal(byteValue, &mockRequestTemplate)

	return &mockRequestTemplate

}

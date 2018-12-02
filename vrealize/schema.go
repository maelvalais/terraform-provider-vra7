package vrealize

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceMachine - use to set resource fields
func ResourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: createResource,
		Read:   readResource,
		Update: updateResource,
		Delete: deleteResource,
		Schema: resourceSchema(),
	}
}

//ResourceSchema - This function is used to update the catalog item template/blueprint
//and replace the values with user defined values added in .tf file.
func ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"catalog_name":        {Type: schema.TypeString, Optional: true},
		"catalog_id":          {Type: schema.TypeString, Optional: true, Computed: true},
		"business_group_id":   {Type: schema.TypeString, Optional: true, Computed: true},
		"business_group_name": {Type: schema.TypeString, Optional: true, Computed: true},
		"wait_timeout":        {Type: schema.TypeInt, Optional: true, Default: 15},
		"request_status":      {Type: schema.TypeString, Optional: false, Computed: true, ForceNew: true},
		"failed_message":      {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
		"deployment_configuration": {Type: schema.TypeMap, Optional: true,
			Elem: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString}},
		"resource_configuration": {Type: schema.TypeMap, Optional: true, Computed: true,
			Elem: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString}},
		"catalog_configuration": {Type: schema.TypeMap, Optional: true,
			Elem: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString}},
	}
}

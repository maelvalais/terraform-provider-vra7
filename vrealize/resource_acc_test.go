package vrealize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gock "gopkg.in/h2non/gock.v1"
)

func init() {
	fmt.Println("init")
}

func TestAccMocked_basic(t *testing.T) {
	gock.DisableNetworking()
	gock.Flush()

	// Two calls are needed because (for some reason) terraform makes two NewClient calls.
	gock.New("http://localhost").Post("/identity/api/tokens").Reply(200).BodyString(resp0Token)
	gock.New("http://localhost").Post("/identity/api/tokens").Reply(200).BodyString(resp0Token)
	gock.New("http://localhost").Get("/catalog-service/api/consumer/entitledCatalogItemViews").ParamPresent("page").ParamPresent("limit").Reply(200).BodyString(resp1EntitledCatalogItemViews)
	gock.New("http://localhost").Get("/catalog-service/api/consumer/entitledCatalogItemViews").Reply(200).BodyString(resp1EntitledCatalogItemViews)
	// For some reason, gock will 'override' previous gock.New() calls when the root is the same.
	// Example: if I do gock.New().Get("/abc") and then gock.New().Get("/abc/new"), and that
	// I do an HTTP request on /abc/new, the Get("/abc") will be returned.....
	// My workaround: first give the 'longuest' GET urls.
	gock.New("http://localhost").Get("/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template").Reply(200).BodyString(resp4Template)
	//gock.New("http://localhost").Get("/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests").Reply(200).BodyString(resp4Request)
	gock.New("http://localhost").Post("/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests").Reply(201).BodyString(resp4Request)

	gock.New("http://localhost").Get("/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be").Reply(200).BodyString(resp2CatalogItemResp)

	gock.New("http://localhost").Get("/catalog-service/api/consumer/requests/b2907df7-6c36-4e30-9c62-a21f293b067a").Reply(200).BodyString(resp5RequestStatus)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true, // disables the need of TF_ACC=1 to enable this test
		PreCheck:   func() { testFieldIsInteger() },
		Providers: map[string]terraform.ResourceProvider{
			"vra7": Provider(),
		},
		CheckDestroy: testNoop(),
		Steps: []resource.TestStep{
			{
				Config: `
				  provider  "vra7" {
						username = "username"
						password  = "password1234!"
						tenant = "vsphere.local"
						host = "http://localhost"
				  }
				  resource "vra7_resource" "resource_1" {
						count = 1
						catalog_name = "CentOS 6.3 - IPAM EXT"
						resource_configuration {
							CentOS_6.3.cpu = "2"
            }
            refresh_seconds = 1
				  }`,
				Check: resource.ComposeTestCheckFunc(
					testFieldIsInteger(),
					resource.TestCheckResourceAttr("resource_configuration", "CentOS_6.3.cpu", "3"),
				),
			},
		},
	})
}

func testNoop() resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return fmt.Errorf("error")

	}
}

func testFieldIsInteger() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		s.Validate()
		field := "CentOS_6.3.cpu"
		if reflect.ValueOf(resourceConfiguration[field]).Kind() != reflect.Int {
			return fmt.Errorf("Field '%s' was expected to be an integer", field)
		}
		return nil
	}
}

var resp0Token = `{
  "expires": "2017-07-25T15:18:49.000Z",
  "id":      "MTUwMDk2NzEyOTEyOTplYTliNTA3YTg4MjZmZjU1YTIwZjp0ZW5hbnQ6dnNwaGVyZS5sb2NhbHVzZXJuYW1lOmphc29uQGNvcnAubG9jYWxleHBpcmF0aW9uOjE1MDA5OTU5MjkwMDA6ZjE1OTQyM2Y1NjQ2YzgyZjY4Yjg1NGFjMGNkNWVlMTNkNDhlZTljNjY3ZTg4MzA1MDViMTU4Y2U3MzBkYjQ5NmQ5MmZhZWM1MWYzYTg1ZWM4ZDhkYmFhMzY3YTlmNDExZmM2MTRmNjk5MGQ1YjRmZjBhYjgxMWM0OGQ3ZGVmNmY=",
  "tenant":  "vsphere.local"
}`

var resp1EntitledCatalogItemViews = `
{
  "links": [
    {
      "@type": "link",
      "rel": "next",
      "href": "https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItemViews?page=2&limit=20"
    }
  ],
  "content": [
    {
      "@type": "ConsumerEntitledCatalogItemView",
      "entitledOrganizations": [
        {
          "tenantRef": "vsphere.local",
          "tenantLabel": "vsphere.local",
          "subtenantRef": "53619006-56bb-4788-9723-9eab79752cc1",
          "subtenantLabel": "Content"
        }
      ],
      "catalogItemId": "e5dd4fba-45ed-4943-b1fc-7f96239286be",
      "name": "CentOS 6.3 - IPAM EXT",
      "description": "CentOS 6.3 IaaS Blueprint w/Infoblox IPAM",
      "isNoteworthy": false,
      "dateCreated": "2016-09-26T13:42:51.564Z",
      "lastUpdatedDate": "2017-01-06T05:11:51.682Z",
      "links": [
        {
          "@type": "link",
          "rel": "GET: Request Template",
          "href": "https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template{?businessGroupId,requestedFor}"
        },
        {
          "@type": "link",
          "rel": "POST: Submit Request",
          "href": "https://vra-01a.corp.local/catalog-service/api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests{?businessGroupId,requestedFor}"
        }
      ],
      "iconId": "e5dd4fba-45ed-4943-b1fc-7f96239286be",
      "catalogItemTypeRef": {
        "id": "com.vmware.csp.component.cafe.composition.blueprint",
        "label": "Composite Blueprint"
      },
      "serviceRef": {
        "id": "baad0ad2-8b96-4347-b188-f534dad53a0d",
        "label": "Infrastructure"
      },
      "outputResourceTypeRef": {
        "id": "composition.resource.type.deployment",
        "label": "Deployment"
      }
    }
  ],
  "metadata": {
    "size": 20,
    "totalElements": 44,
    "totalPages": 3,
    "number": 1,
    "offset": 0
  }
}
`

var resp2CatalogItemResp = `
{
  "catalogItem": {
    "callbacks": null,
    "catalogItemTypeRef": {
      "id": "com.vmware.csp.component.cafe.composition.blueprint",
      "label": "Composite Blueprint"
    },
    "dateCreated": "2015-12-22T03:16:19.289Z",
    "description": "CentOS 6.3 IaaS Blueprint",
    "forms": {
      "itemDetails": {
        "type": "external",
        "formId": "composition.catalog.item.details"
      },
      "catalogRequestInfoHidden": true,
      "requestFormScale": "BIG",
      "requestSubmission": {
        "type": "extension",
        "extensionId": "com.vmware.vcac.core.design.blueprints.requestForm",
        "extensionPointId": null
      },
      "requestDetails": {
        "type": "extension",
        "extensionId": "com.vmware.vcac.core.design.blueprints.requestDetailsForm",
        "extensionPointId": null
      },
      "requestPreApproval": null,
      "requestPostApproval": null
    },
    "iconId": "e5dd4fba-45ed-4943-b1fc-7f96239286be",
    "id": "e5dd4fba-45ed-4943-b1fc-7f96239286be",
    "isNoteworthy": false,
    "lastUpdatedDate": "2017-01-06T05:12:56.690Z",
    "name": "CentOS 6.3",
    "organization": {
      "tenantRef": "vsphere.local",
      "tenantLabel": "vsphere.local",
      "subtenantRef": null,
      "subtenantLabel": null
    },
    "outputResourceTypeRef": {
      "id": "composition.resource.type.deployment",
      "label": "Deployment"
    },
    "providerBinding": {
      "bindingId": "vsphere.local!::!CentOS63",
      "providerRef": {
        "id": "2fbaabc5-3a48-488a-9f2a-a42616345445",
        "label": "Blueprint Service"
      }
    },
    "serviceRef": {
      "id": "baad0ad2-8b96-4347-b188-f534dad53a0d",
      "label": "Infrastructure"
    },
    "status": "PUBLISHED",
    "statusName": "Published",
    "quota": 0,
    "version": 4,
    "requestable": true
  },
  "entitledOrganizations": [
    {
      "tenantRef": "vsphere.local",
      "tenantLabel": "vsphere.local",
      "subtenantRef": "53619006-56bb-4788-9723-9eab79752cc1",
      "subtenantLabel": "Content"
    }
  ]
}

`

var resp4Template = `
{
  "type": "com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest",
  "catalogItemId": "e5dd4fba-45ed-4943-b1fc-7f96239286be",
  "requestedFor": "jason@corp.local",
  "businessGroupId": "53619006-56bb-4788-9723-9eab79752cc1",
  "description": null,
  "reasons": null,
  "data": {
    "CentOS_6.3": {
      "componentTypeId": "com.vmware.csp.component.cafe.composition",
      "componentId": null,
      "classId": "Blueprint.Component.Declaration",
      "typeFilter": "CentOS63*CentOS_6.3",
      "data": {
        "_allocation": {
          "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
          "componentId": null,
          "classId": "Infrastructure.Compute.Machine.Allocation",
          "typeFilter": null,
          "data": {
            "machines": [
              {
                "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                "componentId": null,
                "classId": "Infrastructure.Compute.Machine.Allocation.Machine",
                "typeFilter": null,
                "data": {
                  "machine_id": "",
                  "nics": [
                    {
                      "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                      "componentId": null,
                      "classId": "Infrastructure.Compute.Machine.Nic",
                      "typeFilter": null,
                      "data": {
                        "address": "",
                        "assignment_type": "Static",
                        "external_address": "",
                        "id": null,
                        "load_balancing": null,
                        "network": null,
                        "network_profile": null
                      }
                    }
                  ]
                }
              }
            ]
          }
        },
        "_cluster": 1,
        "_hasChildren": false,
        "cpu": 1,
        "datacenter_location": null,
        "description": "Basic IaaS CentOS Machine",
        "disks": [
          {
            "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
            "componentId": null,
            "classId": "Infrastructure.Compute.Machine.MachineDisk",
            "typeFilter": null,
            "data": {
              "capacity": 3,
              "custom_properties": null,
              "id": 1450725224066,
              "initial_location": "",
              "is_clone": true,
              "label": "Hard disk 1",
              "storage_reservation_policy": "",
              "userCreated": false,
              "volumeId": 0
            }
          }
        ],
        "guest_customization_specification": "CentOS",
        "max_network_adapters": -1,
        "max_per_user": 0,
        "max_volumes": 60,
        "memory": 512,
        "nics": null,
        "os_arch": "x86_64",
        "os_distribution": null,
        "os_type": "Linux",
        "os_version": null,
        "property_groups": null,
        "reservation_policy": null,
        "security_groups": [],
        "security_tags": [],
        "storage": 3
      }
    },
    "_archiveDays": 5,
    "_leaseDays": null,
    "_number_of_instances": 1,
    "corp192168110024": {
      "componentTypeId": "com.vmware.csp.component.cafe.composition",
      "componentId": null,
      "classId": "Blueprint.Component.Declaration",
      "typeFilter": "CentOS63*corp192168110024",
      "data": { "_hasChildren": false }
    }
  }
}
`

var resp4Request = `
{
  "@type": "CatalogItemRequest",
  "id": "b2907df7-6c36-4e30-9c62-a21f293b067a",
  "iconId": "composition.blueprint.png",
  "version": 0,
  "requestNumber": null,
  "state": "PENDING",
  "description": null,
  "reasons": null,
  "requestedFor": "jason@corp.local",
  "requestedBy": "jason@corp.local",
  "organization": {
    "tenantRef": "vsphere.local",
    "tenantLabel": null,
    "subtenantRef": "29a02ed9-7e63-4c77-8a15-c930afb0e3d8",
    "subtenantLabel": null
  },
  "requestorEntitlementId": "e0d6ce92-6e23-4f75-a787-4564699b2895",
  "preApprovalId": null,
  "postApprovalId": null,
  "dateCreated": "2017-08-10T13:38:25.395Z",
  "lastUpdated": "2017-08-10T13:38:25.395Z",
  "dateSubmitted": "2017-08-10T13:38:25.395Z",
  "dateApproved": null,
  "dateCompleted": null,
  "quote": { "leasePeriod": null, "leaseRate": null, "totalLeaseCost": null },
  "requestCompletion": null,
  "requestData": {
    "entries": [
      {
        "key": "MySQL_1",
        "value": {
          "type": "complex",
          "componentTypeId": "com.vmware.csp.component.cafe.composition",
          "componentId": null,
          "classId": "Blueprint.Component.Declaration",
          "typeFilter": "checkcloudclient*MySQL_1",
          "values": {
            "entries": [
              {
                "key": "_hasChildren",
                "value": { "type": "boolean", "value": false }
              },
              {
                "key": "dbpassword",
                "value": {
                  "type": "secureString",
                  "value": "catalog~+gzbqycW+GiAqOREkOs7+mW9D4Og83AKc4FE46i2Z6Y="
                }
              }
            ]
          }
        }
      },
      {
        "key": "Apache_Load_Balancer_1",
        "value": {
          "type": "complex",
          "componentTypeId": "com.vmware.csp.component.cafe.composition",
          "componentId": null,
          "classId": "Blueprint.Component.Declaration",
          "typeFilter": "checkcloudclient*Apache_Load_Balancer_1",
          "values": {
            "entries": [
              {
                "key": "http_node_ips",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "STRING",
                  "items": [{ "type": "string", "value": "None" }]
                }
              },
              {
                "key": "_hasChildren",
                "value": { "type": "boolean", "value": false }
              },
              {
                "key": "http_proxy_port",
                "value": { "type": "string", "value": "8081" }
              },
              { "key": "tomcat_context", "value": null },
              {
                "key": "JAVA_HOME",
                "value": { "type": "string", "value": "/opt/vmware-jre" }
              },
              {
                "key": "appsrv_routes",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "STRING",
                  "items": [{ "type": "string", "value": "None" }]
                }
              },
              {
                "key": "use_ajp",
                "value": { "type": "string", "value": "NO" }
              },
              {
                "key": "http_node_port",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "STRING",
                  "items": [{ "type": "string", "value": "8080" }]
                }
              },
              {
                "key": "http_port",
                "value": { "type": "string", "value": "80" }
              },
              {
                "key": "autogen_sticky_cookie",
                "value": { "type": "string", "value": "NO" }
              }
            ]
          }
        }
      },
      {
        "key": "corp192168110024",
        "value": {
          "type": "complex",
          "componentTypeId": "com.vmware.csp.component.cafe.composition",
          "componentId": null,
          "classId": "Blueprint.Component.Declaration",
          "typeFilter": "checkcloudclient*corp192168110024",
          "values": {
            "entries": [
              {
                "key": "_hasChildren",
                "value": { "type": "boolean", "value": false }
              }
            ]
          }
        }
      },
      {
        "key": "providerId",
        "value": {
          "type": "string",
          "value": "2fbaabc5-3a48-488a-9f2a-a42616345445"
        }
      },
      {
        "key": "subtenantId",
        "value": {
          "type": "string",
          "value": "29a02ed9-7e63-4c77-8a15-c930afb0e3d8"
        }
      },
      {
        "key": "vSphere__vCenter__Machine_2",
        "value": {
          "type": "complex",
          "componentTypeId": "com.vmware.csp.component.cafe.composition",
          "componentId": null,
          "classId": "Blueprint.Component.Declaration",
          "typeFilter": "checkcloudclient*vSphere__vCenter__Machine_2",
          "values": {
            "entries": [
              { "key": "snapshot_name", "value": null },
              { "key": "source_machine", "value": null },
              { "key": "memory", "value": { "type": "integer", "value": 512 } },
              {
                "key": "disks",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "COMPLEX",
                  "items": [
                    {
                      "type": "complex",
                      "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                      "componentId": null,
                      "classId": "Infrastructure.Compute.Machine.MachineDisk",
                      "typeFilter": null,
                      "values": {
                        "entries": [
                          {
                            "key": "is_clone",
                            "value": { "type": "boolean", "value": false }
                          },
                          {
                            "key": "initial_location",
                            "value": { "type": "string", "value": "" }
                          },
                          {
                            "key": "volumeId",
                            "value": { "type": "string", "value": "0" }
                          },
                          {
                            "key": "id",
                            "value": {
                              "type": "integer",
                              "value": 1502347498478
                            }
                          },
                          {
                            "key": "label",
                            "value": { "type": "string", "value": "" }
                          },
                          {
                            "key": "userCreated",
                            "value": { "type": "boolean", "value": true }
                          },
                          {
                            "key": "storage_reservation_policy",
                            "value": { "type": "string", "value": "" }
                          },
                          {
                            "key": "capacity",
                            "value": { "type": "integer", "value": 1 }
                          }
                        ]
                      }
                    }
                  ]
                }
              },
              { "key": "description", "value": null },
              { "key": "storage", "value": { "type": "integer", "value": 1 } },
              { "key": "source_machine_name", "value": null },
              { "key": "guest_customization_specification", "value": null },
              {
                "key": "_hasChildren",
                "value": { "type": "boolean", "value": true }
              },
              { "key": "os_distribution", "value": null },
              { "key": "reservation_policy", "value": null },
              {
                "key": "max_network_adapters",
                "value": { "type": "integer", "value": -1 }
              },
              { "key": "machine_prefix", "value": null },
              {
                "key": "max_per_user",
                "value": { "type": "integer", "value": 0 }
              },
              { "key": "nics", "value": null },
              { "key": "source_machine_vmsnapshot", "value": null },
              {
                "key": "_allocation",
                "value": {
                  "type": "complex",
                  "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                  "componentId": null,
                  "classId": "Infrastructure.Compute.Machine.Allocation",
                  "typeFilter": null,
                  "values": {
                    "entries": [{ "key": "machines", "value": null }]
                  }
                }
              },
              {
                "key": "display_location",
                "value": { "type": "boolean", "value": false }
              },
              { "key": "os_version", "value": null },
              {
                "key": "os_arch",
                "value": { "type": "string", "value": "x86_64" }
              },
              { "key": "cpu", "value": { "type": "integer", "value": 1 } },
              { "key": "datacenter_location", "value": null },
              { "key": "property_groups", "value": null },
              { "key": "_cluster", "value": { "type": "integer", "value": 1 } },
              {
                "key": "security_groups",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "ENTITY_REFERENCE",
                  "items": []
                }
              },
              {
                "key": "max_volumes",
                "value": { "type": "integer", "value": 60 }
              },
              {
                "key": "os_type",
                "value": { "type": "string", "value": "Linux" }
              },
              { "key": "source_machine_external_snapshot", "value": null },
              {
                "key": "security_tags",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "ENTITY_REFERENCE",
                  "items": []
                }
              }
            ]
          }
        }
      },
      { "key": "_leaseDays", "value": null },
      {
        "key": "providerBindingId",
        "value": { "type": "string", "value": "checkcloudclient" }
      },
      {
        "key": "_number_of_instances",
        "value": { "type": "integer", "value": 1 }
      },
      {
        "key": "vSphere__vCenter__Machine_1",
        "value": {
          "type": "complex",
          "componentTypeId": "com.vmware.csp.component.cafe.composition",
          "componentId": null,
          "classId": "Blueprint.Component.Declaration",
          "typeFilter": "checkcloudclient*vSphere__vCenter__Machine_1",
          "values": {
            "entries": [
              { "key": "snapshot_name", "value": null },
              { "key": "source_machine", "value": null },
              { "key": "memory", "value": { "type": "integer", "value": 512 } },
              {
                "key": "disks",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "COMPLEX",
                  "items": [
                    {
                      "type": "complex",
                      "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                      "componentId": null,
                      "classId": "Infrastructure.Compute.Machine.MachineDisk",
                      "typeFilter": null,
                      "values": {
                        "entries": [
                          {
                            "key": "is_clone",
                            "value": { "type": "boolean", "value": false }
                          },
                          {
                            "key": "initial_location",
                            "value": { "type": "string", "value": "hd-1" }
                          },
                          {
                            "key": "volumeId",
                            "value": { "type": "string", "value": "0" }
                          },
                          {
                            "key": "id",
                            "value": {
                              "type": "integer",
                              "value": 1502345335122
                            }
                          },
                          {
                            "key": "label",
                            "value": { "type": "string", "value": "" }
                          },
                          {
                            "key": "userCreated",
                            "value": { "type": "boolean", "value": true }
                          },
                          {
                            "key": "storage_reservation_policy",
                            "value": { "type": "string", "value": "" }
                          },
                          {
                            "key": "capacity",
                            "value": { "type": "integer", "value": 3 }
                          }
                        ]
                      }
                    }
                  ]
                }
              },
              { "key": "description", "value": null },
              { "key": "storage", "value": { "type": "integer", "value": 3 } },
              { "key": "source_machine_name", "value": null },
              { "key": "guest_customization_specification", "value": null },
              {
                "key": "_hasChildren",
                "value": { "type": "boolean", "value": true }
              },
              { "key": "os_distribution", "value": null },
              { "key": "reservation_policy", "value": null },
              {
                "key": "max_network_adapters",
                "value": { "type": "integer", "value": -1 }
              },
              { "key": "machine_prefix", "value": null },
              {
                "key": "max_per_user",
                "value": { "type": "integer", "value": 0 }
              },
              { "key": "nics", "value": null },
              { "key": "source_machine_vmsnapshot", "value": null },
              {
                "key": "_allocation",
                "value": {
                  "type": "complex",
                  "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
                  "componentId": null,
                  "classId": "Infrastructure.Compute.Machine.Allocation",
                  "typeFilter": null,
                  "values": {
                    "entries": [{ "key": "machines", "value": null }]
                  }
                }
              },
              {
                "key": "display_location",
                "value": { "type": "boolean", "value": false }
              },
              { "key": "os_version", "value": null },
              {
                "key": "os_arch",
                "value": { "type": "string", "value": "x86_64" }
              },
              { "key": "cpu", "value": { "type": "integer", "value": 1 } },
              { "key": "datacenter_location", "value": null },
              { "key": "property_groups", "value": null },
              { "key": "_cluster", "value": { "type": "integer", "value": 1 } },
              {
                "key": "security_groups",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "ENTITY_REFERENCE",
                  "items": []
                }
              },
              {
                "key": "max_volumes",
                "value": { "type": "integer", "value": 60 }
              },
              {
                "key": "os_type",
                "value": { "type": "string", "value": "Linux" }
              },
              { "key": "source_machine_external_snapshot", "value": null },
              {
                "key": "security_tags",
                "value": {
                  "type": "multiple",
                  "elementTypeId": "ENTITY_REFERENCE",
                  "items": []
                }
              }
            ]
          }
        }
      }
    ]
  },
  "retriesRemaining": 3,
  "requestedItemName": "myCompositeBlueprint",
  "requestedItemDescription": "",
  "components": null,
  "stateName": null,
  "catalogItemRef": {
    "id": "a3647254-3c50-4fe6-a630-69ae28bf3c81",
    "label": "myCompositeBlueprint"
  },
  "catalogItemProviderBinding": {
    "bindingId": "vsphere.local!::!checkcloudclient",
    "providerRef": {
      "id": "2fbaabc5-3a48-488a-9f2a-a42616345445",
      "label": "Blueprint Service"
    }
  },
  "waitingStatus": "NOT_WAITING",
  "executionStatus": "STARTED",
  "approvalStatus": "PENDING",
  "phase": "PENDING_PRE_APPROVAL"
}
`

var resp5RequestStatus = `
{
  "links": [],
  "content": [
    {
      "@type": "CatalogResourceView",
      "resourceId": "b313acd6-0738-439c-b601-e3ebf9ebb49b",
      "iconId": "502efc1b-d5ce-4ef9-99ee-d4e2a741747c",
      "name": "CentOS 6.3 - IPAM EXT-95563173",
      "description": "",
      "status": null,
      "catalogItemId": "502efc1b-d5ce-4ef9-99ee-d4e2a741747c",
      "catalogItemLabel": "CentOS 6.3 - IPAM EXT",
      "requestId": "dcb12203-93f4-4873-a7d5-1757f3696141",
      "requestState": "SUCCESSFUL",
      "resourceType": "composition.resource.type.deployment",
      "owners": ["Jason Cloud Admin"],
      "businessGroupId": "53619006-56bb-4788-9723-9eab79752cc1",
      "tenantId": "vsphere.local",
      "dateCreated": "2017-07-17T13:26:42.102Z",
      "lastUpdated": "2017-07-17T13:33:25.521Z",
      "lease": { "start": "2017-07-17T13:26:42.079Z", "end": null },
      "costs": null,
      "costToDate": null,
      "totalCost": null,
      "parentResourceId": null,
      "hasChildren": true,
      "data": {},
      "links": [
        {
          "@type": "link",
          "rel": "GET: Catalog Item",
          "href": "http://localhost/catalog-service/api/consumer/entitledCatalogItemViews/502efc1b-d5ce-4ef9-99ee-d4e2a741747c"
        },
        {
          "@type": "link",
          "rel": "GET: Request",
          "href": "http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changelease.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/561be422-ece6-4316-8acb-a8f3dbb8ed0c/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.changeowner.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/59249166-e427-4082-a3dc-eb7223bb2de1/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/3da0ca14-e7e2-4d7b-89cb-c6db57440d72/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.archive.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/9725d56e-461a-471a-be00-b1856681c6d0/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scalein.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/85e090f9-9529-4101-9691-6bab1b0a1f77/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.scaleout.name}",
          "href": "http://localhost/catalog-service/api/consumer/resources/b313acd6-0738-439c-b601-e3ebf9ebb49b/actions/ab5795f5-32ad-4f6c-8598-1d3a7d190caa/requests"
        },
        {
          "@type": "link",
          "rel": "GET: Child Resources",
          "href": "http://localhost/catalog-service/api/consumer/resourceViews?managedOnly=false&withExtendedData=true&withOperations=true&%24filter=parentResource%20eq%20%27b313acd6-0738-439c-b601-e3ebf9ebb49b%27"
        }
      ]
    },
    {
      "@type": "CatalogResourceView",
      "resourceId": "51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5",
      "iconId": "Infrastructure.CatalogItem.Machine.Virtual.vSphere",
      "name": "Content0061",
      "description": "Basic IaaS CentOS Machine",
      "status": "Missing",
      "catalogItemId": null,
      "catalogItemLabel": null,
      "requestId": "dcb12203-93f4-4873-a7d5-1757f3696141",
      "requestState": "SUCCESSFUL",
      "resourceType": "Infrastructure.Virtual",
      "owners": ["Jason Cloud Admin"],
      "businessGroupId": "53619006-56bb-4788-9723-9eab79752cc1",
      "tenantId": "vsphere.local",
      "dateCreated": "2017-07-17T13:33:16.686Z",
      "lastUpdated": "2017-07-17T13:33:25.521Z",
      "lease": { "start": "2017-07-17T13:26:42.079Z", "end": null },
      "costs": null,
      "costToDate": null,
      "totalCost": null,
      "parentResourceId": "b313acd6-0738-439c-b601-e3ebf9ebb49b",
      "hasChildren": false,
      "data": {
        "Component": "CentOS_6.3",
        "DISK_VOLUMES": [
          {
            "componentTypeId": "com.vmware.csp.component.iaas.proxy.provider",
            "componentId": null,
            "classId": "dynamicops.api.model.DiskInputModel",
            "typeFilter": null,
            "data": {
              "DISK_CAPACITY": 3,
              "DISK_INPUT_ID": "DISK_INPUT_ID1",
              "DISK_LABEL": "Hard disk 1"
            }
          }
        ],
        "Destroy": true,
        "EXTERNAL_REFERENCE_ID": "vm-773",
        "IS_COMPONENT_MACHINE": false,
        "MachineBlueprintName": "CentOS 6.3 - IPAM EXT",
        "MachineCPU": 1,
        "MachineDailyCost": 0,
        "MachineDestructionDate": null,
        "MachineExpirationDate": null,
        "MachineGroupName": "Content",
        "MachineGuestOperatingSystem": "CentOS 4/5/6/7 (64-bit)",
        "MachineInterfaceDisplayName": "vSphere (vCenter)",
        "MachineInterfaceType": "vSphere",
        "MachineMemory": 512,
        "MachineName": "Content0061",
        "MachineReservationName": "IPAM Sandbox",
        "MachineStorage": 3,
        "MachineType": "Virtual",
        "NETWORK_LIST": [
          {
            "componentTypeId": "com.vmware.csp.component.iaas.proxy.provider",
            "componentId": null,
            "classId": "dynamicops.api.model.NetworkViewModel",
            "typeFilter": null,
            "data": {
              "NETWORK_ADDRESS": "192.168.110.150",
              "NETWORK_MAC_ADDRESS": "00:50:56:ae:31:bd",
              "NETWORK_NAME": "VM Network",
              "NETWORK_NETWORK_NAME": "ipamext1921681100",
              "NETWORK_PROFILE": "ipam-ext-192.168.110.0"
            }
          }
        ],
        "SNAPSHOT_LIST": [],
        "Unregister": true,
        "VirtualMachine.Admin.UUID": "502e9fb3-6f0d-0b1e-f90f-a769fd406620",
        "endpointExternalReferenceId": "d322b019-58d4-4d6f-9f8b-d28695a716c0",
        "ip_address": "192.168.110.150",
        "machineId": "4fc33663-992d-49f8-af17-df7ce4831aa0"
      },
      "links": [
        {
          "@type": "link",
          "rel": "GET: Request",
          "href": "http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"
        },
        {
          "@type": "link",
          "rel": "GET: Parent Resource",
          "href": "http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}",
          "href": "http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.virtual.Destroy}",
          "href": "http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/654b4c71-e84f-40c7-9439-fd409fea7323/requests"
        },
        {
          "@type": "link",
          "rel": "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}",
          "href": "http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests/template"
        },
        {
          "@type": "link",
          "rel": "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.Unregister}",
          "href": "http://localhost/catalog-service/api/consumer/resources/51bf8bd7-8553-4b0d-b580-41ab0cfaf9a5/actions/f3ae9408-885a-4a3a-9200-43366f2aa163/requests"
        }
      ]
    },
    {
      "@type": "CatalogResourceView",
      "resourceId": "169b596f-e4c0-4b25-ba44-18cb19c0fd65",
      "iconId": "existing_network",
      "name": "ipamext1921681100",
      "description": "Infoblox External Network",
      "status": null,
      "catalogItemId": null,
      "catalogItemLabel": null,
      "requestId": "dcb12203-93f4-4873-a7d5-1757f3696141",
      "requestState": "SUCCESSFUL",
      "resourceType": "Infrastructure.Network.Network.Existing",
      "owners": ["Jason Cloud Admin"],
      "businessGroupId": "53619006-56bb-4788-9723-9eab79752cc1",
      "tenantId": "vsphere.local",
      "dateCreated": "2017-07-17T13:27:17.526Z",
      "lastUpdated": "2017-07-17T13:33:25.521Z",
      "lease": { "start": "2017-07-17T13:26:42.079Z", "end": null },
      "costs": null,
      "costToDate": null,
      "totalCost": null,
      "parentResourceId": "b313acd6-0738-439c-b601-e3ebf9ebb49b",
      "hasChildren": false,
      "data": {
        "Description": "Infoblox External Network",
        "IPAMEndpointId": "1c2b6237-540a-43c3-8c06-b37a1d274b44",
        "IPAMEndpointName": "Infoblox - nios01a",
        "Name": "ipamext1921681100",
        "_archiveDays": 5,
        "_hasChildren": false,
        "_leaseDays": null,
        "_number_of_instances": 1,
        "dns": {
          "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
          "componentId": null,
          "classId": "Infrastructure.Network.Network.DnsWins",
          "typeFilter": null,
          "data": {
            "alternate_wins": null,
            "dns_search_suffix": null,
            "dns_suffix": null,
            "preferred_wins": null,
            "primary_dns": null,
            "secondary_dns": null
          }
        },
        "gateway": null,
        "ip_ranges": [
          {
            "componentTypeId": "com.vmware.csp.iaas.blueprint.service",
            "componentId": null,
            "classId": "Infrastructure.Network.Network.IpRanges",
            "typeFilter": null,
            "data": {
              "description": "",
              "end_ip": "",
              "externalId": "network/default-vra/192.168.110.0/24",
              "id": "b078d23a-1c3d-4458-ab57-e352c80e6d55",
              "name": "192.168.110.0/24",
              "start_ip": ""
            }
          }
        ],
        "network_profile": "ipam-ext-192.168.110.0",
        "providerBindingId": "CentOS63Infoblox",
        "providerId": "2fbaabc5-3a48-488a-9f2a-a42616345445",
        "subnet_mask": "255.255.255.0",
        "subtenantId": "53619006-56bb-4788-9723-9eab79752cc1"
      },
      "links": [
        {
          "@type": "link",
          "rel": "GET: Request",
          "href": "http://localhost/catalog-service/api/consumer/requests/dcb12203-93f4-4873-a7d5-1757f3696141"
        },
        {
          "@type": "link",
          "rel": "GET: Parent Resource",
          "href": "http://localhost/catalog-service/api/consumer/resourceViews/b313acd6-0738-439c-b601-e3ebf9ebb49b"
        }
      ]
    }
  ],
  "metadata": {
    "size": 20,
    "totalElements": 3,
    "totalPages": 1,
    "number": 1,
    "offset": 0
  }
}
`

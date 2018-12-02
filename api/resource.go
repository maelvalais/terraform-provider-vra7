package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vmware/terraform-provider-vra7/utils"
)

//ResourceActionTemplate - is used to store information
//related to resource action template information.
type ResourceActionTemplate struct {
	Type        string                 `json:"type,omitempty"`
	ResourceID  string                 `json:"resourceId,omitempty"`
	ActionID    string                 `json:"actionId,omitempty"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

//ResourceView - is used to store information
//related to resource template information.
type ResourceView struct {
	Content []interface {
	} `json:"content"`
	Links []interface{} `json:"links"`
}

//RequestStatusView - used to store REST response of
//request triggered against any resource.
type RequestStatusView struct {
	RequestCompletion struct {
		RequestCompletionState string `json:"requestCompletionState"`
		CompletionDetails      string `json:"CompletionDetails"`
	} `json:"requestCompletion"`
	Phase string `json:"phase"`
}

type BusinessGroups struct {
	Content []BusinessGroup `json:"content,omitempty"`
}

type BusinessGroup struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id,omitempty"`
}

// Resource View of a provisioned request
type RequestResourceView struct {
	Content []DeploymentResource `json:"content,omitempty"`
	Links   []interface{}        `json:"links,omitempty"`
}

type DeploymentResource struct {
	RequestState    string                 `json:"requestState,omitempty"`
	Description     string                 `json:"description,omitempty"`
	LastUpdated     string                 `json:"lastUpdated,omitempty"`
	TenantId        string                 `json:"tenantId,omitempty"`
	Name            string                 `json:"name,omitempty"`
	BusinessGroupId string                 `json:"businessGroupId,omitempty"`
	DateCreated     string                 `json:"dateCreated,omitempty"`
	Status          string                 `json:"status,omitempty"`
	RequestId       string                 `json:"requestId,omitempty"`
	ResourceId      string                 `json:"resourceId,omitempty"`
	ResourceType    string                 `json:"resourceType,omitempty"`
	ResourcesData   DeploymentResourceData `json:"data,omitempty"`
}

type DeploymentResourceData struct {
	Memory                      int    `json:"MachineMemory,omitempty"`
	Cpu                         int    `json:"MachineCPU,omitempty"`
	IpAddress                   string `json:"ip_address,omitempty"`
	Storage                     int    `json:"MachineStorage,omitempty"`
	MachineInterfaceType        string `json:"MachineInterfaceType,omitempty"`
	MachineName                 string `json:"MachineName,omitempty"`
	MachineGuestOperatingSystem string `json:"MachineGuestOperatingSystem,omitempty"`
	MachineDestructionDate      string `json:"MachineDestructionDate,omitempty"`
	MachineGroupName            string `json:"MachineGroupName,omitempty"`
	MachineBlueprintName        string `json:"MachineBlueprintName,omitempty"`
	MachineReservationName      string `json:"MachineReservationName,omitempty"`
	MachineType                 string `json:"MachineType,omitempty"`
	MachineId                   string `json:"machineId,omitempty"`
	MachineExpirationDate       string `json:"MachineExpirationDate,omitempty"`
	Component                   string `json:"Component,omitempty"`
	Expire                      bool   `json:"Expire,omitempty"`
	Reconfigure                 bool   `json:"Reconfigure,omitempty"`
	Reset                       bool   `json:"Reset,omitempty"`
	Reboot                      bool   `json:"Reboot,omitempty"`
	PowerOff                    bool   `json:"PowerOff,omitempty"`
	Destroy                     bool   `json:"Destroy,omitempty"`
	Shutdown                    bool   `json:"Shutdown,omitempty"`
	Suspend                     bool   `json:"Suspend,omitempty"`
	Reprovision                 bool   `json:"Reprovision,omitempty"`
	ChangeLease                 bool   `json:"ChangeLease,omitempty"`
	ChangeOwner                 bool   `json:"ChangeOwner,omitempty"`
	CreateSnapshot              bool   `json:"CreateSnapshot,omitempty"`
}

// Retrieves the resources that were provisioned as a result of a given request.
// Also returns the actions allowed on the resources and their templates
type ResourceActions struct {
	Links   []interface{}           `json:"links,omitempty"`
	Content []ResourceActionContent `json:"content,omitempty"`
}

type ResourceActionContent struct {
	Id              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	ResourceTypeRef ResourceTypeRef `json:"resourceTypeRef,omitempty"`
	Status          string          `json:"status,omitempty"`
	RequestId       string          `json:"requestId,omitempty"`
	RequestState    string          `json:"requestState,omitempty"`
	Operations      []Operation     `json:"operations,omitempty"`
	ResourceData    ResourceDataMap `json:"resourceData,omitempty"`
}

type ResourceTypeRef struct {
	Id    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type Operation struct {
	Name        string `json:"name,omitempty"`
	OperationId string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type ResourceDataMap struct {
	Entries []ResourceDataEntry `json:"entries,omitempty"`
}
type ResourceDataEntry struct {
	Key   string                 `json:"key,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

//CatalogRequest - A structure that captures a vRA catalog request.
type CatalogRequest struct {
	ID           string      `json:"id"`
	IconID       string      `json:"iconId"`
	Version      int         `json:"version"`
	State        string      `json:"state"`
	Description  string      `json:"description"`
	Reasons      interface{} `json:"reasons"`
	RequestedFor string      `json:"requestedFor"`
	RequestedBy  string      `json:"requestedBy"`
	Organization struct {
		TenantRef      string `json:"tenantRef"`
		TenantLabel    string `json:"tenantLabel"`
		SubtenantRef   string `json:"subtenantRef"`
		SubtenantLabel string `json:"subtenantLabel"`
	} `json:"organization"`

	RequestorEntitlementID   string                 `json:"requestorEntitlementId"`
	PreApprovalID            string                 `json:"preApprovalId"`
	PostApprovalID           string                 `json:"postApprovalId"`
	DateCreated              time.Time              `json:"dateCreated"`
	LastUpdated              time.Time              `json:"lastUpdated"`
	DateSubmitted            time.Time              `json:"dateSubmitted"`
	DateApproved             time.Time              `json:"dateApproved"`
	DateCompleted            time.Time              `json:"dateCompleted"`
	Quote                    interface{}            `json:"quote"`
	RequestData              map[string]interface{} `json:"requestData"`
	RequestCompletion        string                 `json:"requestCompletion"`
	RetriesRemaining         int                    `json:"retriesRemaining"`
	RequestedItemName        string                 `json:"requestedItemName"`
	RequestedItemDescription string                 `json:"requestedItemDescription"`
	Components               string                 `json:"components"`
	StateName                string                 `json:"stateName"`

	CatalogItemProviderBinding struct {
		BindingID   string `json:"bindingId"`
		ProviderRef struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"providerRef"`
	} `json:"catalogItemProviderBinding"`

	Phase           string `json:"phase"`
	ApprovalStatus  string `json:"approvalStatus"`
	ExecutionStatus string `json:"executionStatus"`
	WaitingStatus   string `json:"waitingStatus"`
	CatalogItemRef  struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"catalogItemRef"`
}

//GetBusinessGroupId retrieves business group id from business group name
func (vRAClient *Client) GetBusinessGroupId(businessGroupName string) (string, error) {

	path := "/identity/api/tenants/" + vRAClient.Tenant + "/subtenants?%24filter=name+eq+'" + businessGroupName + "'"
	log.Info("Fetching business group id from name..GET %s ", path)
	BusinessGroups := new(BusinessGroups)
	apiError := new(Error)
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(BusinessGroups, apiError)
	if err != nil {
		return "", err
	}
	if !apiError.IsEmpty() {
		return "", apiError
	}
	// BusinessGroups array will contain only one BusinessGroup element containing the BG
	// with the name businessGroupName.
	// Fetch the id of that BG
	return BusinessGroups.Content[0].Id, nil
}

//DestroyMachine - To set resource destroy call
func (vRAClient *Client) DestroyMachine(destroyTemplate *ResourceActionTemplate, destroyActionURL string) error {
	apiError := new(Error)
	resp, err := vRAClient.HTTPClient.New().Post(destroyActionURL).
		BodyJSON(destroyTemplate).Receive(nil, apiError)

	if resp.StatusCode != 201 {
		log.Errorf("The destroy deployment request failed with error: %v ", resp.Status)
		return err
	}

	if !apiError.IsEmpty() {
		log.Errorf("The destroy deployment request failed with error: %v ", apiError.Errors)
		return apiError
	}
	return nil
}

//GetRequestStatus - To read request status of resource
// which is used to show information to user post create call.
func (vRAClient *Client) GetRequestStatus(requestId string) (*RequestStatusView, error) {
	//Form a URL to read request status
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s", requestId)
	RequestStatusViewTemplate := new(RequestStatusView)
	apiError := new(Error)
	//Set a REST call and fetch a resource request status
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(RequestStatusViewTemplate, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.IsEmpty() {
		return nil, apiError
	}
	return RequestStatusViewTemplate, nil
}

// GetDeploymentState - Read the state of a vRA7 Deployment
func (vRAClient *Client) GetDeploymentState(CatalogRequestId string) (*ResourceView, error) {
	//Form an URL to fetch resource list view
	path := fmt.Sprintf("catalog-service/api/consumer/requests/%s"+
		"/resourceViews", CatalogRequestId)
	ResourceView := new(ResourceView)
	apiError := new(Error)
	//Set a REST call to fetch resource view data
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(ResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.IsEmpty() {
		return nil, apiError
	}
	return ResourceView, nil
}

// Retrieves the resources that were provisioned as a result of a given request.
func (vRAClient *Client) GetRequestResourceView(catalogRequestId string) (*RequestResourceView, error) {
	path := fmt.Sprintf(utils.GET_REQUEST_RESOURCE_VIEW_API, catalogRequestId)
	requestResourceView := new(RequestResourceView)
	apiError := new(Error)
	_, err := vRAClient.HTTPClient.New().Get(path).Receive(requestResourceView, apiError)
	if err != nil {
		return nil, err
	}
	if !apiError.IsEmpty() {
		return nil, apiError
	}
	return requestResourceView, nil
}

//RequestCatalogItem - Make a catalog request.
func (vRAClient *Client) RequestCatalogItem(requestTemplate *CatalogItemRequestTemplate) (*CatalogRequest, error) {
	//Form a path to set a REST call to create a machine
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s"+
		"/requests", requestTemplate.CatalogItemID)

	catalogRequest := new(CatalogRequest)
	apiError := new(Error)

	jsonBody, jErr := json.Marshal(requestTemplate)
	if jErr != nil {
		log.Error("Error marshalling request templat as JSON")
		return nil, jErr
	}

	log.Info("JSON Request Info: %s", string(jsonBody))
	//Set a REST call to create a machine
	_, err := vRAClient.HTTPClient.New().Post(path).BodyJSON(requestTemplate).
		Receive(catalogRequest, apiError)

	if err != nil {
		return nil, err
	}

	if !apiError.IsEmpty() {
		return nil, apiError
	}

	return catalogRequest, nil
}

package api

import (
	"errors"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func MockClient() Client {
	return NewClient(
		"admin@myvra.local",
		"pass!@#",
		"vsphere.local",
		"http://localhost/",
		true,
	)
}

func TestClient_Authenticate(t *testing.T) {
	client := MockClient()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(200, `{
		  "expires": "2017-07-25T15:18:49.000Z",
		  "id": "MTUwMDk2NzEyOTEyOTplYTliNTA3YTg4MjZmZjU1YTIwZjp0ZW5hbnQ6dnNwaGVyZS5sb2NhbHVzZXJuYW1lOmphc29uQGNvcnAubG9jYWxleHBpcmF0aW9uOjE1MDA5OTU5MjkwMDA6ZjE1OTQyM2Y1NjQ2YzgyZjY4Yjg1NGFjMGNkNWVlMTNkNDhlZTljNjY3ZTg4MzA1MDViMTU4Y2U3MzBkYjQ5NmQ5MmZhZWM1MWYzYTg1ZWM4ZDhkYmFhMzY3YTlmNDExZmM2MTRmNjk5MGQ1YjRmZjBhYjgxMWM0OGQ3ZGVmNmY=",
		  "tenant": "vsphere.local"
		}`))

	err := client.Authenticate()

	if len(client.BearerToken) == 0 {
		t.Error("Fail to set BearerToken.")
	}

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewErrorResponder(errors.New(`{
		  "errors": [
			{
			  "code": 90135,
			  "source": null,
			  "message": "Unable to authenticate user jason@corp.local1 in tenant vsphere.local.",
			  "systemMessage": "90135-Unable to authenticate user jason@corp.local1 in tenant vsphere.local.",
			  "moreInfoUrl": null
			}
		  ]
		}`)))

	err = client.Authenticate()

	if err == nil {
		t.Errorf("Authentication should fail")
	}
}

func TestNewClient(t *testing.T) {
	username := "admin@myvra.local"
	password := "pass!@#"
	tenant := "vshpere.local"
	baseURL := "http://localhost/"

	client := NewClient(
		username,
		password,
		tenant,
		baseURL,
		true,
	)

	if client.Username != username {
		t.Errorf("Expected username %v, got %v ", username, client.Username)
	}

	if client.Password != password {
		t.Errorf("Expected password %v, got %v ", password, client.Password)
	}

	if client.Tenant != tenant {
		t.Errorf("Expected tenant %v, got %v ", tenant, client.Tenant)
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseUrl %v, got %v ", baseURL, client.BaseURL)
	}
}

# terraform-provider-vra7-brmc

My goal with this fork is to _adapt_ the upstream project
[terraform-provider-vra7] to the non-conventional ways BRMC is implementing XaaS
blueprints. As-is, [terraform-provider-vra7] cannot work for multiple reasons:

- Right now, I get an error after 15 minutes after posting the request:

      PROVIDER_FAILED (Workflow:Wait for resource action request / Extract result (item3)#11)

  No idea why, this error also happens when creating a VM from the web client.

- One field in the returned catalog item template is missing, which makes it
  hard to guess that it is needed when the error messages returned by vRA are so
  cryptic, e.g., for `bgName` missing, it was:

      No domain name. Look for data binding in form (Workflow:Provision a VM / check values (item27)#1)

- The BRMC team has chosen a weird layout for catalog item templates. In the
  upstream [terraform-provided-vra], properties (cpu, codeBasicat...) are
  expected to be wrapped into components, not directly properties in the
  top-level "data" field. Here is what a template looks like in BRMC:

  ```json
  {
    "type": "com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest",
    "catalogItemId": "37d038c8-7a7c-41dd-9beb-0420c08fc815",
    "requestedFor": "NLGN2101@ad.francetelecom.fr",
    "businessGroupId": "21e6cf47-952d-4b67-9881-439af6388a41",
    "description": null,
    "reasons": null,
    "data": {
      "cpu": null,
      "codeBasicat": null,
      "cactusGroupNames": null,
      "module_applicatif": null,
      "AZ": null
    }
  }
  ```

  Instead, they should use components, e.g., `Platon_rehl73`:

  ```json
  {
      "type": "com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest",
      "catalogItemId": "37d038c8-7a7c-41dd-9beb-0420c08fc815",
      "requestedFor": "NLGN2101@ad.francetelecom.fr",
      "businessGroupId": "21e6cf47-952d-4b67-9881-439af6388a41",
      "description": null,
      "reasons": null,

      "data": {
          "Platon_rehl73": {
              "componentTypeId": "...",
              "componentId": null,
              "classId": "Blueprint.Component.Declaration",
              "typeFilter": "Platon_rehl73*Platon_rehl73",
              "data": {
                  "Cafe.Shim.VirtualMachine.MaxCost": 0,
                  "Cafe.Shim.VirtualMachine.MinCost": 0,
                  "_cluster": 1,
                  "_hasChildren": false,
                  "action": "FullClone",
                  "allow_storage_policies": false,
                  "archive_days": 0,
                  "blueprint_type": "1",
                  "cpu": 1,
                  "custom_properties": [
                      "codeBasicat": null,
                      "cactusGroupNames": null,
                      "module_applicatif": null,
                      "AZ": null
                  ]
          ...
      }
  }
  ```

  (Note: this example comes from the [vra-75-programming-guide], page 51)

  This components/properties thing is explained in the upstream [terraform-provider-vra7] project
  in [comments](https://github.com/vmware/terraform-provider-vra7/blob/88688609cd8d848c17cb124646f2d90709741c47/vrealize/resource.go#L109-L112):

  ```go
  // User-supplied resource configuration keys are expected to be of the form:
  //     <component name>.<property name>.
  // Extract the property names and values for each component in the blueprint, and add/update
  // them in the right location in the request template.
  ```

- Unhelpful/cryptic error messages when provisionning fails or when the API call is wrong.
  Try to guess what this one means:

      PROVIDER_FAILED (Workflow:Wait for resource action request / Extract result (item3)#11)

[terraform-provider-vra7]: https://github.com/vmware/terraform-provider-vra7
[vra-75-programming-guide]: https://vdc-download.vmware.com/vmwb-repository/dcr-public/3f9ab622-499d-4caf-801b-2a6c1f83a6d4/ba2d7e2c-3320-4dfb-9ef0-39588cadda2e/vrealize-automation-75-programming-guide.pdf

A self-contained deployable integration between Terraform and vRealize Automation (vRA) which allows Terraform users to request/provision entitled vRA catalog items using Terraform. Supports Terraform destroying vRA provisioned resources.

## Getting Started

These instructions will get you a copy of the project up and run on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Prerequisites

To get the vRA plugin up and running you need the following things.

- [Terraform 0.9 or above](https://www.terraform.io/downloads.html)
- [Go Language 1.9.2 or above](https://golang.org/dl/)
- [dep - new dependency management tool for Go](https://github.com/golang/dep)

## Project Setup

Setup a GOLang project structure

```
|-/home/<USER>/TerraformPluginProject
    |-bin
    |-pkg
    |-src

```

## Environment Setup

Set following environment variables

**Linux Users**

_GOROOT is a golang library path_

```
export GOROOT=/usr/local/go
```

_GOPATH is a path pointing toward the source code directory_

```
export GOPATH=/home/<user>/TerraformPluginProject
```

**Windows Users**

_GOROOT is a golang library path_

```
set GOROOT=C:\Go
```

_GOPATH is a path pointing toward the source code directory_

```
set GOPATH=C:\TerraformPluginProject
```

## Set terraform provider

**Linux Users**

Create _~/.terraformrc_ and put following content in it.

```
    providers {
         vra7 = "/home/<USER>/TerraformPluginProject/bin/terraform-provider-vra7"
    }
```

**Windows Users**

Create _%APPDATA%/terraform.rc_ and put following content in it.

```
    providers {
         vra7 = "C:\\TerraformPluginProject\\bin\\terraform-provider-vra7.exe"
    }
```

## Installation

Clone repo code into go project using _go get_

```
    go get github.com/vmware/terraform-provider-vra7

```

## Create Binary

**Linux and MacOS Users**

Navigate to _/home/<USER>/TerraformPluginProject/src/github.com/vmware/terraform-provider-vra7_ and run go build command to generate plugin binary

```
    dep ensure
    go build -o /home/<USER>/TerraformPluginProject/bin/terraform-provider-vra7

```

**Windows Users**

Navigate to _C:\TerraformPluginProject\src\github.com\vmware\terraform-provider-vra7_ and run go build command to generate plugin binary

```
    dep ensure
    go build -o C:\TerraformPluginProject\bin\terraform-provider-vra7.exe

```

## Create Terraform Configuration file

The VMware vRA terraform configuration file contains two objects

### Provider

This part contains service provider details.

**Configure Provider**

Provider block contains four mandatory fields

- **username** - _vRA portal username_
- **password** - _vRA portal password_
- **tenant** - _vRA portal tenant_
- **host** - _End point of REST API_
- **insecure** - _In case of self-signed certificates. Default value is false._

Example

```
    provider "vra7" {
      username = "vRAUser1@vsphere.local"
      password = "password123!"
      tenant = "corp.local.tenant"
      host = "http://myvra.example.com/"
      insecure = false
    }

```

### Resource

This part contains any resource that can be deployed on that service provider.
For example, in our case machine blueprint, software blueprint, complex blueprint, network, etc.

**Configure Resource**

Syntax

```
resource "vra7" "<resource_name1>" {
}
```

The resource block contains mandatory and optional fields as follows:

Mandatory:

One of catalog_name or catalog_id must be specified in the resource configuration.

- **catalog_name** - _catalog_name is a field which contains valid catalog name from your vRA_

- **catalog_id** - _catalog_id is a field which contains a valid catalog id from your vRA._

Optional:

- **businessgroup_id** - _This is an optional field. You can specify a different Business Group ID from what provided by default in the template reques, provided that your account is allowed to do it_

- **catalog_configuration** - _This is an optional field. If catalog properties have default values or no mandatory user input required for catalog service then you can skip this field from the terraform configuration file. This field contains user inputs to catalog services. Value of this field is a key value pair. Key is any field name of catalog and value is any valid user input to the respective field._

- **count** - _This field is used to create replicas of resources. If count is not provided then it will be considered as 1 by default._

- **deployment_configuration** - _This is an optional field. Can only be used to specify the description or reasons field at the deployment level. Key is any field name of catalog and value is any valid user input to the respective field._

- **resource_configuration** - _This is an optional field. If blueprint properties have default values or no mandatory property value is required then you can skip this field from terraform configuration file. This field contains user inputs to catalog services. Value of this field is in key value pair. Key is service.field_name and value is any valid user input to the respective field._

- **wait_timeout** - _This is an optional field with a default value of 15. It defines the time to wait (in minutes) for a resource operation to complete successfully._

Example 1

```
resource "vra7_resource" "example_machine1" {
  catalog_name = "CentOS 6.3"
   resource_configuration = {
         Linux.cpu = "1"
         Windows2008R2SP1.cpu =  "2"
         Windows2012.cpu =  "4"
         Windows2016.cpu =  "2"
     }
     catalog_configuration = {
         lease_days = "5"
     }
     deployment_configuration = {
         reasons      = "I have some"
         description  = "deployment via terraform"
     }
     count = 3
}

```

Example 2

```
resource "vra7_resource" "example_machine2" {
  catalog_id = "e5dd4fba7f96239286be45ed"
   resource_configuration = {
         Linux.cpu = "1"
         Windows2008.cpu =  "2"
         Windows2012.cpu =  "4"
         Windows2016.cpu =  "2"
     }
     count = 4
}

```

Save this configuration in main.tf in a path where the binary is placed.

## Execution

These are the Terraform commands that can be used for the vRA plugin:

- **terraform init** - _The init command is used to initialize a working directory containing Terraform configuration files._

- **terraform plan** - _Plan command shows plan for resources like how many resources will be provisioned and how many will be destroyed._

- **terraform apply** - _apply is responsible to execute actual calls to provision resources._

- **terraform refresh** - _By using the refresh command you can check the status of the request._

- **terraform show** - _show will set a console output for resource configuration and request status._

- **terraform destroy** - _destroy command will destroy all the resources present in terraform configuration file._

Navigate to the location where main.tf and binary are placed and use the above commands as needed.

## Contributing

The terraform-provider-vra7 project team welcomes contributions from the community. Before you start working with terraform-provider-vra7, please read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License

terraform-provider-vra7 is available under the [MIT license](LICENSE).

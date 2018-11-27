provider "vra7" {
  username = "%s"
  password = "%s"
  tenant   = "vsphere.local"
  host     = "https://brmc.si.fr.intraorange"
  insecure = true
}

resource "vra7_resource" "vm" {
  count = 1

  catalog_name = "Provisionner une VM DCaaS"

  //businessgroup_name = "DOM-4"
  businessgroup_id = "21e6cf47-952d-4b67-9881-439af6388a41"

  # catalog_name = "Provisionner une VM FAST"

  resource_configuration = {
    vmSuffix                   = "dvadxws00bxxxxx"
    typeServeur                = "dv"                                            # dv (custom naming)
    serverType                 = "Développement"
    typeServerFullName         = "Développement"
    predefinedRole             = "ws00"
    role                       = "ws00"
    customRole                 = true
    customRoleValue            = "ws00"
    module_applicatif          = "PF CC AGILE DELIVERY - DEV - DESI"
    OS                         = "Pl@ton Linux RedHat"
    region                     = "Normandie"
    AZ                         = "Salle 4"
    securityGroupName          = "SGIC-DOM-4-super-flux-N1"
    cpu                        = 1
    ram                        = 1
    diskData                   = 0
    targetDiskSizeOfVm         = 20
    currentDiskSizeOfBlueprint = 20
    cos                        = "Standard"
    leaseUnlimited             = false
    lease                      = 90
    niveauSupport              = "P2"
    addDRSGroup                = false
    backupPlanned              = "Désactivée"
    codeBasicat                = "ADX"                                           // Example: ERB
    commentaire                = "petite VM de test"                             // Example: A VM for testing
    cos                        = "Standard"
    labelThresholdCPU          = "CPU limit crossed - you will need approval"
    labelCPU                   = "1"                                             // Ne marche pas avec '3 vCPU' par exemple; je pense qu'ils parsent ce label
    labelThresholdMemory       = "Memory limit crossed - you will need approval"
    labelRam                   = "1GB"                                           // On peut y mettre n'importe quoi
    currentDiskSizeOfBlueprint = 20
    diskData                   = 0
    labelDataDiskSize          = 0
    targetDiskSizeOfVm         = 20
    domainType                 = "DCaaS"
    drsGroupDesc               = ""
    groupeDRS                  = ""
    hasBasicat                 = true
    isBGFast                   = false
    niveauSupport              = "P2"                                            // P2 = lowest support 8*5 a week, P1 = higher support, 24*7 a week
    region                     = "Normandie"
    searchBasicat              = "adx"
    securityGroupName          = "SGIC-DOM-4-super-flux-N1"                      // example: 
    useCloudInit               = false
    cloudInitData              = ""                                              // "${data.template_file.cloud_init.rendered}"
    supportEntity              = "/Orange/Of/Dtsi/Desi/Dixsi/Ptal/Pre"
    clientEntity               = "/Orange/Of/Dtsi/Desi/Dixsi/Ptal/Pre"
    bgName                     = "DOM-4"

    # cactusGroupNames           = ""                                  // Example: [BRM000011] (see in cactus)
  }
  refresh_seconds = 10 // seconds
  wait_timeout    = 30 // minutes
  catalog_configuration = {
    reasons     = "Test"
    description = "deployment via terraform"
  }
}

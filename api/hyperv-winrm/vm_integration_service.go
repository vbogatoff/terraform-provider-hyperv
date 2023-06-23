package hyperv_winrm

import (
	"context"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type getVmIntegrationServicesArgs struct {
	VmName string
}

var getVmIntegrationServicesTemplate = template.Must(template.New("GetVmIntegrationServices").Parse(`
$ErrorActionPreference = 'Stop'
$vmIntegrationServicesObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMIntegrationService | %{ @{
	Name=$_.Name;
	Enabled=$_.Enabled;
}})

if ($vmIntegrationServicesObject) {
	$vmIntegrationServices = ConvertTo-Json -InputObject $vmIntegrationServicesObject
	$vmIntegrationServices
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmIntegrationServices(ctx context.Context, vmName string) (result []api.VmIntegrationService, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVmIntegrationServicesTemplate, getVmIntegrationServicesArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type enableVmIntegrationServiceArgs struct {
	VmName string
	Name   string
}

var enableVmIntegrationServiceTemplate = template.Must(template.New("EnableVmIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

integrationServiceId := ""
	switch '{{.Name}}' {
	case "Time Synchronization":
		integrationServiceId = "2497F4DE-E9FA-4204-80E4-4B75C46419C0"
	case "Heartbeat":
		integrationServiceId = "84EAAE65-2F2E-45F5-9BB5-0E857DC8EB47"
	case "Key-Value Pair Exchange":
		integrationServiceId = "2A34B1C2-FD73-4043-8A5B-DD2159BC743F"
	case "Shutdown":
		integrationServiceId = "9F8233AC-BE49-4C79-8EE3-E7E1985B2077"
	case "VSS":
		integrationServiceId = "5CED1297-4598-4915-A5FC-AD21BB4D02A4"
	case "Guest Service Interface":
		integrationServiceId = "6C09BB55-D683-4DA0-8931-C9BF705F6480"
	default:
		panic("unrecognized Integration Service Name")

Get-VMIntegrationService -VmName '{{.VmName}}' | ?{$_.Id -match $integrationServiceId} | Enable-VMIntegrationService

`))

func (c *ClientConfig) EnableVmIntegrationService(ctx context.Context, vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, enableVmIntegrationServiceTemplate, enableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

type disableVmIntegrationServiceArgs struct {
	VmName string
	Name   string
}

var disableVmIntegrationServiceTemplate = template.Must(template.New("DisableVmIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

integrationServiceId := ""
	switch '{{.Name}}' {
	case "Time Synchronization":
		integrationServiceId = "2497F4DE-E9FA-4204-80E4-4B75C46419C0"
	case "Heartbeat":
		integrationServiceId = "84EAAE65-2F2E-45F5-9BB5-0E857DC8EB47"
	case "Key-Value Pair Exchange":
		integrationServiceId = "2A34B1C2-FD73-4043-8A5B-DD2159BC743F"
	case "Shutdown":
		integrationServiceId = "9F8233AC-BE49-4C79-8EE3-E7E1985B2077"
	case "VSS":
		integrationServiceId = "5CED1297-4598-4915-A5FC-AD21BB4D02A4"
	case "Guest Service Interface":
		integrationServiceId = "6C09BB55-D683-4DA0-8931-C9BF705F6480"
	default:
		panic("unrecognized Integration Service Name")

Get-VMIntegrationService -VmName '{{.VmName}}' | ?{$_.Id -match $integrationServiceId} | Disable-VMIntegrationService

`))

func (c *ClientConfig) DisableVmIntegrationService(ctx context.Context, vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, disableVmIntegrationServiceTemplate, disableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmIntegrationServices(ctx context.Context, vmName string, integrationServices []api.VmIntegrationService) (err error) {
	for _, integrationService := range integrationServices {
		if integrationService.Enabled {
			err = c.EnableVmIntegrationService(ctx, vmName, integrationService.Name)
		} else {
			err = c.DisableVmIntegrationService(ctx, vmName, integrationService.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

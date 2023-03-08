package azureFacade

import "testing"

func TestAutoChangeValidIpForVM(t *testing.T) {
	azureParams := AzureParams{
		AzureClientId:        "xxxx",
		AzureClientSecret:    "xxxx",
		AzureTenantId:        "xxxx",
		SubscriptionId:       "xxxx",
		ResourceGroupName:    "xxxx",
		Location:             "Japan East",
		NetworkInterfaceName: "xxxx",
		IpConfigurationName:  "", // can be ""
	}
	err := AutoChangeValidIpForVM(azureParams, 3)
	if err != nil {
		t.Fatal(err.Error())
	}
}

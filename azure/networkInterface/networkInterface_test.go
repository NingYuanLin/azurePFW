package networkInterface

import (
	"github.com/NingYuanLin/azurePFW/azure/auth"
	"context"
	"testing"
)

func TestUpdatePublicIpId(t *testing.T) {
	var (
		azureClientId        = "xxx"
		azureClientSecret    = "xxx"
		azureTenantId        = "xxx"
		subscriptionId       = "xxx"
		resourceGroupName    = "xxx"
		location             = "Japan East"
		networkInterfaceName = "xxx"
		ipConfigurationName  = "" // can be ""
	)
	publicIpId := "xxx"  // such as "/subscriptions/xxx/resourceGroups/i_group/providers/Microsoft.Network/publicIPAddresses/xxx"

	cred, err := auth.GetCred(azureClientId, azureClientSecret, azureTenantId)
	if err != nil {
		t.Fatalf("auth err:%+v\n", err)
	}

	client, err := NewClient(cred, subscriptionId, resourceGroupName, location, networkInterfaceName, ipConfigurationName)
	nic, err := client.UpdatePublicIpId(context.Background(), publicIpId)
	if err != nil {
		t.Fatalf("CreateOrUpdateNIC err:%+v\n", err)
	}
	_ = nic
}

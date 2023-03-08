package publicIpAddress

import (
	"github.com/NingYuanLin/azurePFW/azure/auth"
	"context"
	"fmt"
	"testing"
)

func TestCreateOrUpdatePublicIp(t *testing.T) {
	t.Log("Be careful: The new created ip not be deleted automatically. So you should delete it manually")
	var (
		azureClientId       = "xxx"
		azureClientSecret   = "xxx"
		azureTenantId       = "xxx"
		subscriptionId      = "xxx"
		resourceGroupName   = "xxx"
		location            = "Japan East"
		publicIpAddressName = "xxx"
	)
	cred, err := auth.GetCred(azureClientId, azureClientSecret, azureTenantId)
	if err != nil {
		t.Fatalf("auth err:%+v\n", err)
	}

	client, err := NewClient(cred, subscriptionId, resourceGroupName, location)
	if err != nil {
		t.Fatalf("client create err:%+v", err)
	}

	ip, err := client.CreateOrUpdatePublicIp(context.Background(), publicIpAddressName)
	if err != nil {
		t.Fatalf("create ip err: %+v", err)
	}
	t.Logf("ip.ID: %s\n", *ip.ID)
}

func TestDeletePublicIp(t *testing.T) {
	var (
		azureClientId       = "xxx"
		azureClientSecret   = "xxx"
		azureTenantId       = "xxx"
		subscriptionId      = "xxx"
		resourceGroupName   = "xxx"
		location            = "Japan East"
	)
	publicIPAddressName := "xxx"

	cred, err := auth.GetCred(azureClientId, azureClientSecret, azureTenantId)
	if err != nil {
		t.Fatalf("auth err:%+v\n", err)
	}

	client, err := NewClient(cred, subscriptionId, resourceGroupName, location)
	if err != nil {
		t.Fatalf("client create err:%+v", err)
	}
	err = client.DeletePublicIp(context.Background(), publicIPAddressName)
	if err != nil {
		t.Fatalf("delete public ip err:%+v\n", err)
	}
}

func TestGetSpecificIpInfo(t *testing.T) {
	var (
		azureClientId       = "xxx"
		azureClientSecret   = "xxx"
		azureTenantId       = "xxx"
		subscriptionId      = "xxx"
		resourceGroupName   = "xxx"
		location            = "Japan East"
	)
	publicIPAddressName := "xxx"

	cred, err := auth.GetCred(azureClientId, azureClientSecret, azureTenantId)
	if err != nil {
		t.Fatalf("auth err:%+v\n", err)
	}

	client, err := NewClient(cred, subscriptionId, resourceGroupName, location)
	if err != nil {
		t.Fatalf("client create err:%+v", err)
	}
	ipInfo, err := client.GetSpecificIpInfo(context.Background(), publicIPAddressName)
	if err != nil {
		t.Fatalf("GetSpecificIpInfo err:%s\n", err.Error())
	}
	fmt.Println(*ipInfo.Properties.IPAddress)
}

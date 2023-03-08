package networkInterface

import (
	"azurePFW/utils/logger"
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

type Client struct {
	SubscriptionId       string
	ResourceGroupName    string
	Location             string
	Cred                 *azidentity.DefaultAzureCredential
	NetworkInterfaceName string
	IpConfigurationName  string // when equal to "", it means automatic detection
	NicClient            *armnetwork.InterfacesClient
}

func NewClient(cred *azidentity.DefaultAzureCredential, subscriptionId, resourceGroupName, location, networkInterfaceName, ipConfigurationName string) (*Client, error) {
	c := &Client{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		Location:             location,
		NetworkInterfaceName: networkInterfaceName,
		IpConfigurationName:  ipConfigurationName,
	}
	c.Cred = cred
	nicClient, err := armnetwork.NewInterfacesClient(c.SubscriptionId, c.Cred, nil)
	if err != nil {
		return nil, err
	}
	c.NicClient = nicClient
	return c, nil
}

func (c *Client) CreateOrUpdateNIC(ctx context.Context, request armnetwork.Interface) (*armnetwork.Interface, error) {
	pollerResp, err := c.NicClient.BeginCreateOrUpdate(
		ctx,
		c.ResourceGroupName,
		c.NetworkInterfaceName,
		request,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Interface, nil
}

func (c *Client) GetNicInfo(ctx context.Context) (*armnetwork.InterfacesClientGetResponse, error) {
	info, err := c.NicClient.Get(ctx, c.ResourceGroupName, c.NetworkInterfaceName, nil)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *Client) GetSpecificIpConfiguration(ipConfigurations []*armnetwork.InterfaceIPConfiguration) (*armnetwork.InterfaceIPConfiguration, error) {
	var ipConfiguration *armnetwork.InterfaceIPConfiguration
	if len(ipConfigurations) == 0 {
		return nil, errors.New("无ip configuration配置")
	} else {
		if c.IpConfigurationName == "" {
			if len(ipConfigurations) > 1 {
				logger.Info("发现多个ip configuration配置且您未指定ipConfigurationName，默认使用第一个")
			}
			ipConfiguration = ipConfigurations[0]
		} else {
			for _, v := range ipConfigurations {
				if *v.Name == c.IpConfigurationName {
					ipConfiguration = v
					break
				}
			}
			if ipConfiguration == nil {
				return nil, errors.New("无法找到您所指定的ipConfigurationName")
			}
		}
	}
	return ipConfiguration, nil
}

func (c *Client) UpdatePublicIpId(ctx context.Context, publicIpId string) (*armnetwork.Interface, error) {
	info, err := c.GetNicInfo(ctx)
	if err != nil {
		return nil, err
	}

	ipConfigurations := info.Properties.IPConfigurations
	ipConfiguration, err := c.GetSpecificIpConfiguration(ipConfigurations)
	if err != nil {
		return nil, err
	}
	ipConfiguration.Properties.PublicIPAddress.ID = to.Ptr(publicIpId)

	nicResponse, err := c.CreateOrUpdateNIC(ctx, info.Interface)
	if err != nil {
		return nil, err
	}
	return nicResponse, nil
}

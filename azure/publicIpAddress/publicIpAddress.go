package publicIpAddress

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

type Client struct {
	SubscriptionId    string
	ResourceGroupName string
	Location          string
	Cred              *azidentity.DefaultAzureCredential
	PublicIpClient    *armnetwork.PublicIPAddressesClient
}

func NewClient(cred *azidentity.DefaultAzureCredential, subscriptionId, resourceGroupName, location string) (*Client, error) {
	c := &Client{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		Location:          location,
	}

	c.Cred = cred
	publicIpClient, err := armnetwork.NewPublicIPAddressesClient(c.SubscriptionId, c.Cred, nil)
	if err != nil {
		return nil, err
	}
	c.PublicIpClient = publicIpClient
	return c, nil
}

func (c *Client) CreateOrUpdatePublicIp(ctx context.Context, publicIPAddressName string) (*armnetwork.PublicIPAddress, error) {
	pollerResp, err := c.PublicIpClient.BeginCreateOrUpdate(
		ctx,
		c.ResourceGroupName,
		publicIPAddressName,
		armnetwork.PublicIPAddress{
			Name:     to.Ptr(publicIPAddressName),
			Location: to.Ptr(c.Location),
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.PublicIPAddress, nil
}

func (c *Client) DeletePublicIp(ctx context.Context, publicIPAddressName string) error {
	pollerResp, err := c.PublicIpClient.BeginDelete(
		ctx,
		c.ResourceGroupName,
		publicIPAddressName,
		nil,
	)
	if err != nil {
		return err
	}
	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSpecificIpInfo(ctx context.Context, publicIPAddressName string) (armnetwork.PublicIPAddressesClientGetResponse, error) {
	publicIpAddressInfo, err := c.PublicIpClient.Get(ctx, c.ResourceGroupName, publicIPAddressName, nil)
	return publicIpAddressInfo, err
}

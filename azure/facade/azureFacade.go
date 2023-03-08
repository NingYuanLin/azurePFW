package azureFacade

import (
	"azurePFW/azure/auth"
	"azurePFW/azure/networkInterface"
	"azurePFW/azure/publicIpAddress"
	"azurePFW/utils/logger"
	"azurePFW/utils/net/icmp"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type AzureParams struct {
	AzureClientId        string
	AzureClientSecret    string
	AzureTenantId        string
	SubscriptionId       string
	ResourceGroupName    string
	Location             string
	NetworkInterfaceName string
	IpConfigurationName  string // when equal to "", it means automatic detection
}

func AutoChangeValidIpForVM(azureParams AzureParams, retryMaxNum int) error {
	for retryNum := 1; retryNum <= retryMaxNum; retryNum++ {
		cred, err := auth.GetCred(azureParams.AzureClientId, azureParams.AzureClientSecret, azureParams.AzureTenantId)
		if err != nil {
			return err
			//panic(err.Error())
		}

		networkInterfaceClient, err := networkInterface.NewClient(
			cred,
			azureParams.SubscriptionId,
			azureParams.ResourceGroupName,
			azureParams.Location,
			azureParams.NetworkInterfaceName,
			azureParams.IpConfigurationName,
		)
		if err != nil {
			return err
			//panic(err.Error())
		}
		publicIpAddressClient, err := publicIpAddress.NewClient(cred, azureParams.SubscriptionId, azureParams.ResourceGroupName, azureParams.Location)
		if err != nil {
			//panic(err.Error())
			return err
		}

		nicInfo, err := networkInterfaceClient.GetNicInfo(context.Background())
		if err != nil {
			return err
		}
		ipConfiguration, err := networkInterfaceClient.GetSpecificIpConfiguration(nicInfo.Properties.IPConfigurations)
		if err != nil {
			return err
		}
		oldIpAddrId := *ipConfiguration.Properties.PublicIPAddress.ID
		oldIpAddrIdSplit := strings.Split(oldIpAddrId, "/")
		oldIpAddrName := oldIpAddrIdSplit[len(oldIpAddrIdSplit)-1]
		oldIpInfo, err := publicIpAddressClient.GetSpecificIpInfo(context.Background(), oldIpAddrName)
		if err != nil {
			return err
		}
		oldIpAddr := *oldIpInfo.Properties.IPAddress

		canReach := false
		for retryNum := 0; retryNum < 5; retryNum++ {
			canReach = icmp.CanReach(oldIpAddr)
			if canReach == true {
				break
			}
		}

		if canReach == true {
			logger.Info("ip地址有效")
			return nil
		}
		logger.Info("ip地址无效")

		// 创建新ip
		logger.Info("尝试创建新ip")
		newIpAddrName := strconv.FormatInt(time.Now().UnixMilli(), 10)
		publicIpAddrInfo, err := publicIpAddressClient.CreateOrUpdatePublicIp(context.Background(), newIpAddrName)
		if err != nil {
			return err
		}

		// 绑定nic到新ip
		logger.Infof("新ip创建成功, 新ip: %s", *publicIpAddrInfo.Properties.IPAddress)
		logger.Info("尝试绑定nic到新ip")
		_, err = networkInterfaceClient.UpdatePublicIpId(context.Background(), *publicIpAddrInfo.ID)
		if err != nil {
			return err
		}

		// 删除旧ip, 可以与其他操作并行
		logger.Info("尝试删除旧ip")
		deleteOldIpErrChan := make(chan error, 1)
		go func() {
			err = publicIpAddressClient.DeletePublicIp(context.Background(), oldIpAddrName)
			deleteOldIpErrChan <- err
		}()

		// 测试新ip是否可以ping通
		isNewIpOk := false
		for ipTestRetryNum := 0; ipTestRetryNum < 10; ipTestRetryNum++ {
			logger.Infof("retry:%d 测试新ip是否可以ping通，请确保防火墙允许icmp访问", ipTestRetryNum)
			canReach = icmp.CanReach(*publicIpAddrInfo.Properties.IPAddress)
			if canReach == true {
				logger.Info("新ip可以ping通,成功")
				isNewIpOk = true
				break
			}
			time.Sleep(time.Second * 2)
		}

		deleteOldIpErr := <-deleteOldIpErrChan
		if deleteOldIpErr != nil {
			return deleteOldIpErr
		}

		if isNewIpOk == true {
			return nil
		}

		time.Sleep(time.Second)
	}
	return errors.New("超过最大尝试次数，仍无法找到可用ip")
}

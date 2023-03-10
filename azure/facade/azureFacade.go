package azureFacade

import (
	"github.com/NingYuanLin/azurePFW/azure/auth"
	"github.com/NingYuanLin/azurePFW/azure/networkInterface"
	"github.com/NingYuanLin/azurePFW/azure/publicIpAddress"
	"github.com/NingYuanLin/azurePFW/utils/logger"
	"github.com/NingYuanLin/azurePFW/utils/net/icmp"
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
			logger.Info("ip????????????")
			return nil
		}
		logger.Info("ip????????????")

		// ?????????ip
		logger.Info("???????????????ip")
		newIpAddrName := strconv.FormatInt(time.Now().UnixMilli(), 10)
		publicIpAddrInfo, err := publicIpAddressClient.CreateOrUpdatePublicIp(context.Background(), newIpAddrName)
		if err != nil {
			return err
		}

		// ??????nic??????ip
		logger.Infof("???ip????????????, ???ip: %s", *publicIpAddrInfo.Properties.IPAddress)
		logger.Info("????????????nic??????ip")
		_, err = networkInterfaceClient.UpdatePublicIpId(context.Background(), *publicIpAddrInfo.ID)
		if err != nil {
			return err
		}

		// ?????????ip, ???????????????????????????
		logger.Info("???????????????ip")
		deleteOldIpErrChan := make(chan error, 1)
		go func() {
			err = publicIpAddressClient.DeletePublicIp(context.Background(), oldIpAddrName)
			deleteOldIpErrChan <- err
		}()

		// ?????????ip????????????ping???
		isNewIpOk := false
		for ipTestRetryNum := 0; ipTestRetryNum < 10; ipTestRetryNum++ {
			logger.Infof("retry:%d ?????????ip????????????ping??????????????????????????????icmp??????", ipTestRetryNum)
			canReach = icmp.CanReach(*publicIpAddrInfo.Properties.IPAddress)
			if canReach == true {
				logger.Info("???ip??????ping???,??????")
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
	return errors.New("????????????????????????????????????????????????ip")
}

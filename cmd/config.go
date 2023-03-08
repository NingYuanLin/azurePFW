package cmd

import (
	"github.com/NingYuanLin/azurePFW/azure/facade"
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "The operations about config file",
	Run: func(cmd *cobra.Command, args []string) {
		if create, _ := cmd.Flags().GetBool("create"); create == true {
			// create config file
			err := createConfigFile()
			if err != nil {
				log.Println(err)
			}
			return
		}
	},
}

func init() {
	configCmd.Flags().Bool("create", false, "create config file")
	configCmd.MarkFlagRequired("create")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//fmt.Println("initConfig")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		//print(home)
		cobra.CheckErr(err)

		// Search config in home directory with name "azurePFW.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("azurePFW")
	}

	//viper.AutomaticEnv() // read in environment variables that match

}

func readConfigToViper() (configFileUsed string, err error) {
	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	configFileUsed = viper.ConfigFileUsed()
	return
}

func createConfigFile() error {
	// check if the config file exist and read to viper
	configFile, err := readConfigToViper()
	exist := false
	if err == nil {
		exist = true
	}

	if exist {
		//	log.Println("Using config file:", configFile)
		log.Printf("Warning: Config file is existed in %s, and the following operation will rewrite the config file.\n", configFile)
	}

	reader := bufio.NewReader(os.Stdin)
	var azureClientId string
	for {
		fmt.Print("Please input azure client id: ")
		azureClientId, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		azureClientId = strings.TrimSpace(azureClientId)
		if azureClientId != "" {
			break
		}
	}

	var azureClientSecret string
	for {
		fmt.Print("Please input azure client secret: ")
		azureClientSecret, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		azureClientSecret = strings.TrimSpace(azureClientSecret)
		if azureClientSecret != "" {
			break
		}
	}

	var azureTenantId string
	for {
		fmt.Print("Please input azure tenant id: ")
		azureTenantId, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		azureTenantId = strings.TrimSpace(azureTenantId)
		if azureTenantId != "" {
			break
		}
	}

	var subscriptionId string
	for {
		fmt.Print("Please input subscription id: ")
		subscriptionId, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		subscriptionId = strings.TrimSpace(subscriptionId)
		if subscriptionId != "" {
			break
		}
	}

	var resourceGroupName string
	for {
		fmt.Print("Please input resource group name: ")
		resourceGroupName, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		resourceGroupName = strings.TrimSpace(resourceGroupName)
		if resourceGroupName != "" {
			break
		}
	}

	var location string
	for {
		fmt.Print("Please input the location of the resource (such as \"Japan East\"): ")
		location, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		location = strings.TrimSpace(location)
		if location != "" {
			break
		}
	}

	var networkInterfaceName string
	for {
		fmt.Print("Please input the network interface name: ")
		networkInterfaceName, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		networkInterfaceName = strings.TrimSpace(networkInterfaceName)
		if networkInterfaceName != "" {
			break
		}
	}

	var ipConfigurationName string
	fmt.Print("Please input the ip configuration name (it will be detected automatically by default): ")
	ipConfigurationName, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	ipConfigurationName = strings.TrimSpace(ipConfigurationName)

	viper.Set("azure_client_id", azureClientId)
	viper.Set("azure_client_secret", azureClientSecret)
	viper.Set("azure_tenant_id", azureTenantId)
	viper.Set("subscription_id", subscriptionId)
	viper.Set("resource_group_name", resourceGroupName)
	viper.Set("location", location)
	viper.Set("network_interface_name", networkInterfaceName)
	viper.Set("ip_configuration_name", ipConfigurationName)

	if exist == true {
		err = viper.WriteConfig()
	} else {
		err = viper.SafeWriteConfig()
	}
	if err != nil {
		return err
	}

	return nil
}

func parseConfigFromFile() (azureFacade.AzureParams, error) {
	azureParams := azureFacade.AzureParams{}
	configFile, err := readConfigToViper()
	if err != nil {
		return azureParams, err
	}
	log.Println("Using config file:", configFile)

	azureParams.AzureClientId = viper.GetString("azure_client_id")
	azureParams.AzureClientSecret = viper.GetString("azure_client_secret")
	azureParams.AzureTenantId = viper.GetString("azure_tenant_id")
	azureParams.SubscriptionId = viper.GetString("subscription_id")
	azureParams.ResourceGroupName = viper.GetString("resource_group_name")
	azureParams.Location = viper.GetString("location")
	azureParams.NetworkInterfaceName = viper.GetString("network_interface_name")
	azureParams.IpConfigurationName = viper.GetString("ip_configuration_name")

	return azureParams, nil
}

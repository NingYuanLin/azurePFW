package cmd

import (
	"github.com/NingYuanLin/azurePFW/azure/facade"
	"github.com/NingYuanLin/azurePFW/utils/logger"
	"github.com/spf13/cobra"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "To start azurePFW",
	Run: func(cmd *cobra.Command, args []string) {
		azureParams, err := parseConfigFromFile()
		if err != nil {
			panic(err)
		}

		for {
			err = azureFacade.AutoChangeValidIpForVM(azureParams, 3)
			if err != nil {
				logger.Error(err)
			}
			time.Sleep(time.Second * 3)
		}

	},
}

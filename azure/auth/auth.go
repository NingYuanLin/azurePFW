package auth

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"os"
)

func GetCred(azureClientId, azureClientSecret, azureTenantId string) (*azidentity.DefaultAzureCredential, error) {
	os.Setenv("AZURE_CLIENT_ID", azureClientId)
	os.Setenv("AZURE_CLIENT_SECRET", azureClientSecret)
	os.Setenv("AZURE_TENANT_ID", azureTenantId)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//log.Fatalf("Authentication failure: %+v", err)
		return nil, errors.New(fmt.Sprintf("Authentication failure: \n%s", err.Error()))
	}
	return cred, nil
}

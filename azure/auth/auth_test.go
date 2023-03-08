package auth

import "testing"

func TestGetCred(t *testing.T) {
	var (
		azureClientId     = ""
		azureClientSecret = ""
		azureTenantId     = ""
	)
	cred, err := GetCred(azureClientId, azureClientSecret, azureTenantId)
	if err != nil {
		t.Fatalf("err: %+v\n", err)
	}
	_ = cred
}

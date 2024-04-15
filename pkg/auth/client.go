package auth

import (
	"encoding/json"
	"fmt"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
)

var EnvOS = openstack.NewEnv(envPrefix)

const envPrefix = "OS_"

type cc struct {
	*openstack.Cloud
	*golangsdk.ProviderClient
}

func copyCloud(src *openstack.Cloud) (*openstack.Cloud, error) {
	srcJson, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("error marshalling cloud: %s", err)
	}
	res := new(openstack.Cloud)
	if err := json.Unmarshal(srcJson, res); err != nil {
		return nil, fmt.Errorf("error unmarshalling cloud: %s", err)
	}
	return res, nil
}

func CloudAndClient() (*cc, error) {
	cloud, err := EnvOS.Cloud()
	if err != nil {
		return nil, fmt.Errorf("error constructing cloud configuration: %w", err)
	}
	cloud, err = copyCloud(cloud)
	if err != nil {
		return nil, fmt.Errorf("error copying cloud: %w", err)
	}
	client, err := EnvOS.AuthenticatedClient()
	if err != nil {
		return nil, err
	}
	return &cc{cloud, client}, nil
}

func SetupTemporaryAKSK(config *cc) error {
	if config.AKSKAuthOptions.AccessKey != "" {
		return nil
	}

	client, err := NewIdentityV3Client()
	if err != nil {
		return fmt.Errorf("error creating identity v3 domain client: %s", err)
	}
	credential, err := credentials.CreateTemporary(client, credentials.CreateTemporaryOpts{
		Methods: []string{"token"},
		Token:   client.Token(),
	}).Extract()
	if err != nil {
		return fmt.Errorf("error creating temporary AK/SK: %s", err)
	}

	config.AKSKAuthOptions.AccessKey = credential.AccessKey
	config.AKSKAuthOptions.SecretKey = credential.SecretKey
	config.AKSKAuthOptions.SecurityToken = credential.SecurityToken
	return nil
}

func NewIdentityV3Client() (*golangsdk.ServiceClient, error) {
	cc, err := CloudAndClient()
	if err != nil {
		return nil, err
	}

	return openstack.NewIdentityV3(cc.ProviderClient, golangsdk.EndpointOpts{
		Region: cc.RegionName,
	})
}

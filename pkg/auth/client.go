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

// cc stands for `cloud` & `client`
type cc struct {
	*openstack.Cloud
	*golangsdk.ProviderClient
}

// copyCloud makes a deep copy of cloud configuration
func copyCloud(src *openstack.Cloud) (*openstack.Cloud, error) {
	// Marshal the source cloud struct to JSON
	srcJson, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("error marshalling cloud: %s", err)
	}
	// Create a new cloud struct and unmarshal the JSON into it
	res := new(openstack.Cloud)
	if err := json.Unmarshal(srcJson, res); err != nil {
		return nil, fmt.Errorf("error unmarshalling cloud: %s", err)
	}
	return res, nil
}

// CloudAndClient returns a copy of the cloud configuration and an authenticated client for the OS_ environment
func CloudAndClient() (*cc, error) {
	// Retrieve cloud configuration from environment variables
	cloud, err := EnvOS.Cloud()
	if err != nil {
		return nil, fmt.Errorf("error constructing cloud configuration: %w", err)
	}
	// Make a deep copy of the cloud configuration
	cloud, err = copyCloud(cloud)
	if err != nil {
		return nil, fmt.Errorf("error copying cloud: %w", err)
	}
	// Authenticate the client using the cloud configuration
	client, err := EnvOS.AuthenticatedClient()
	if err != nil {
		return nil, err
	}
	return &cc{cloud, client}, nil
}

// SetupTemporaryAKSK configures temporary AK/SK credentials for the given cloud configuration
func SetupTemporaryAKSK(config *cc) error {
	// If AccessKey is already set, no need to create temporary credentials
	if config.AKSKAuthOptions.AccessKey != "" {
		return nil
	}

	// Create a new Identity V3 client
	client, err := NewIdentityV3Client()
	if err != nil {
		return fmt.Errorf("error creating identity v3 domain client: %s", err)
	}

	// Create temporary AK/SK credentials using the client
	credential, err := credentials.CreateTemporary(client, credentials.CreateTemporaryOpts{
		Methods: []string{"token"},
		Token:   client.Token(),
	}).Extract()
	if err != nil {
		return fmt.Errorf("error creating temporary AK/SK: %s", err)
	}

	// Set the temporary credentials in the cloud configuration
	config.AKSKAuthOptions.AccessKey = credential.AccessKey
	config.AKSKAuthOptions.SecretKey = credential.SecretKey
	config.AKSKAuthOptions.SecurityToken = credential.SecurityToken
	return nil
}

// NewIdentityV3Client creates a new Identity V3 client for interacting with the OpenStack Identity service
func NewIdentityV3Client() (*golangsdk.ServiceClient, error) {
	// Retrieve the cloud configuration and authenticated client
	cc, err := CloudAndClient()
	if err != nil {
		return nil, err
	}

	// Create a new Identity V3 service client
	return openstack.NewIdentityV3(cc.ProviderClient, golangsdk.EndpointOpts{
		Region: cc.RegionName,
	})
}

package obs

import (
	"fmt"
	"io"
	"log"

	"github.com/opentelekomcloud-infra/vault-raft-backup/pkg/auth"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/obs"
)

const obsACL = "private"

// Client wraps an OBS client to provide simplified object storage operations.
type Client struct {
	obsClient *obs.ObsClient // obsClient is the underlying client from the OBS service.
}

// Config contains configuration details for OBS operations, such as bucket and object names.
type Config struct {
	BucketName string
	ObjectName string
}

// NewOBSClient initializes a new OBS client using credentials and endpoints fetched via auth package.
func NewOBSClient() (*Client, error) {
	cc, err := auth.CloudAndClient() // Get cloud and client configurations.
	if err != nil {
		return nil, err
	}

	if err := auth.SetupTemporaryAKSK(cc); err != nil { // Set up temporary Access Key and Secret Key.
		return nil, fmt.Errorf("failed to construct OBS client without AK/SK: %s", err)
	}

	client, err := openstack.NewOBSService(cc.ProviderClient, golangsdk.EndpointOpts{
		Region: cc.RegionName,
	})
	if err != nil {
		return nil, err
	}
	opts := cc.AKSKAuthOptions
	obsClient, err := obs.New(
		opts.AccessKey, opts.SecretKey, client.Endpoint,
		obs.WithSecurityToken(opts.SecurityToken), obs.WithSignature(obs.SignatureObs),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OBS client: %s", err)
	}
	return &Client{obsClient: obsClient}, nil
}

// InitClient provides an alias to NewOBSClient for initializing an OBS client.
func InitClient() (*Client, error) {
	return NewOBSClient()
}

// UploadToOBS uploads a file to the OBS using the provided reader, bucket, and object configuration.
func (c *Client) UploadToOBS(buffer io.Reader, config *Config) error {
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: config.BucketName,
				Key:    config.ObjectName,
				ACL:    obsACL,
			},
		},
		Body: buffer,
	}

	output, err := c.obsClient.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to upload to OBS: %s", err)
	}

	log.Printf("Upload to OBS successful: %+v", output)
	return nil
}

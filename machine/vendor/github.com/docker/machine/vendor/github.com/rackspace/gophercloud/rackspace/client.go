package rackspace

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/gophercloud/rackspace/identity/v2/tokens"
)

const (
	v20 = "v2.0"
)

// AuthenticatedClient logs in to Rackspace with the provided credentials and
// constructs a ProviderClient with the returned credentials.
func AuthenticatedClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewClient prepares an unauthenticated ProviderClient for the Rackspace
// endpoint. Most users will prefer using the AuthenticatedClient function
// instead.
//
// See http://docs.rackspace.com/auth/api/v2.0/auth-client-devguide/content/QuickStart-000.html
func NewClient(endpoint string) (*gophercloud.ProviderClient, error) {
	if endpoint == "" {
		endpoint = USIdentityEndpoint
	}

	switch endpoint {
	case USIdentityEndpoint, UKIdentityEndpoint:
		return os.NewClient(endpoint)
	default:
		return os.NewClient(endpoint)
	}
}

// Authenticate or re-authenticate against the most appropriate identity service
// given the provided options.
func Authenticate(client *gophercloud.ProviderClient, options gophercloud.AuthOptions) error {
	v := options.AllowReauth

	return v2auth(client, "", options, gophercloud.EndpointOpts{Availability: gophercloud.AvailabilityPublic})
}

func v2auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions, eopts gophercloud.EndpointOpts) error {
	v2Client, err := os.NewIdentityV2(client, eopts)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens.AuthOptions{
		Options: options,
	}

	result := tokens.Create(v2Client, v2Opts)
	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eopts)
		}
	}
	client.TokenID = token.ID
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return os.V2EndpointURL(catalog, opts)
	}

	if v {
		client.ReauthFunc = func() error {
			return Authenticate(client, options)
		}
	}

	return nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the
// v2 identity service.
func NewIdentityV2(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewIdentityV2(client, eopts)
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute
// package.
func NewComputeV2(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewComputeV2(client, eopts)
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network
// package.
func NewNetworkV2(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewNetworkV2(client, eopts)
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1
// object storage package.
func NewObjectStorageV1(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewObjectStorageV1(client, eopts)
}

// NewBlockStorageV1 creates a ServiceClient that may be used with the v1
// block storage package.
func NewBlockStorageV1(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewBlockStorageV1(client, eopts)
}

// NewCDNV1 creates a ServiceClient that may be used with the v1 CDN package.
func NewCDNV1(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewCDNV1(client, eopts)
}

// NewOrchestrationV1 creates a ServiceClient that may be used with the v1 orchestration
// package.
func NewOrchestrationV1(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewOrchestrationV1(client, eopts)
}

// NewDBV1 creates a ServiceClient that may be used with the v1 DB package.
func NewDBV1(client *gophercloud.ProviderClient, eopts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return os.NewDBV1(client, eopts)
}

// V2EndpointURL returns the endpoint URL for the specified service in the
// Rackspace Service Catalog.
func V2EndpointURL(catalog *ostokens.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	return os.V2EndpointURL(catalog, opts)
}

// Failure is a JSON representation of an error returned from the Rackspace API.
type Failure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (f Failure) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.Message)
}

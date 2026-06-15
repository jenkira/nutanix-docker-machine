package tokens

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
)

// AuthOptions wraps the OpenStack base AuthOptions struct with Rackspace
// extensions, such as APIKey.
type AuthOptions struct {
	os.AuthOptions

	// APIKey is the Rackspace-specific credential for API key-based
	// authentication.
	APIKey string `json:"-"`
}

// Create authenticates to the identity v2 endpoint using the provided
// options.
func Create(c *gophercloud.ServiceClient, auth AuthOptions) os.CreateResult {
	type APIKeyCredentials struct {
		Username string `json:"username"`
		APIKey   string `json:"apiKey"`
	}

	type Auth struct {
		*os.PasswordCredentials `json:"passwordCredentials,omitempty"`
		*APIKeyCredentials      `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
		TenantID                string `json:"tenantId,omitempty"`
		TenantName              string `json:"tenantName,omitempty"`
		Token                   *os.TokenCredentials `json:"token,omitempty"`
	}

	type request struct {
		Auth Auth `json:"auth"`
	}

	reqBody := request{
		Auth: Auth{
			TenantID:   auth.TenantID,
			TenantName: auth.TenantName,
		},
	}

	if auth.APIKey != "" {
		reqBody.Auth.APIKeyCredentials = &APIKeyCredentials{
			Username: auth.Username,
			APIKey:   auth.APIKey,
		}
	} else if auth.TokenID != "" {
		reqBody.Auth.Token = &os.TokenCredentials{ID: auth.TokenID}
	} else {
		reqBody.Auth.PasswordCredentials = &os.PasswordCredentials{
			Username: auth.Username,
			Password: auth.Password,
		}
	}

	var result os.CreateResult
	_, result.Err = c.Post(os.TokensURL(c), reqBody, &result.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return result
}

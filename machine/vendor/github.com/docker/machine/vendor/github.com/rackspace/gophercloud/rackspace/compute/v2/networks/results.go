package networks

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Network represents a network resource.
type Network struct {
	// UUID for the network
	ID string `mapstructure:"id"`

	// Human-readable name for the network. Might not be unique.
	Label string `mapstructure:"label"`

	// CIDR block for the network
	CIDR string `mapstructure:"cidr"`
}

// GetResult is returned from a call to the Get function.
type GetResult struct {
	gophercloud.SingletonResult
}

// Extract interprets GetResults as a Network, if possible.
func (r GetResult) Extract() (*Network, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Network *Network `mapstructure:"network"`
	}

	err := gophercloud.DecodeResponse(r.Body, &response)
	if err != nil {
		return nil, err
	}

	return response.Network, nil
}

// NetworkPage stores a page of Network results.
type NetworkPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a NetworkPage is empty.
func (page NetworkPage) IsEmpty() (bool, error) {
	ns, err := ExtractNetworks(page)
	return len(ns) == 0, err
}

// ExtractNetworks interprets a NetworkPage as a slice of Networks.
func ExtractNetworks(page pagination.Page) ([]Network, error) {
	var response struct {
		Networks []Network `mapstructure:"networks"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.Networks, nil
}

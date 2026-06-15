package virtualinterfaces

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// VirtualInterface represents a virtual interface resource.
type VirtualInterface struct {
	// UUID for the virtual interface
	ID string `mapstructure:"id"`

	// MAC address for the virtual interface
	MACAddress string `mapstructure:"mac_address"`

	// IPAddresses for the virtual interface
	IPAddresses []IPAddress `mapstructure:"ip_addresses"`
}

// IPAddress represents an IP address resource.
type IPAddress struct {
	NetworkID   string `mapstructure:"network_id"`
	NetworkLabel string `mapstructure:"network_label"`
	Address      string `mapstructure:"address"`
}

// VirtualInterfacePage stores a page of VirtualInterface results.
type VirtualInterfacePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a VirtualInterfacePage is empty.
func (page VirtualInterfacePage) IsEmpty() (bool, error) {
	vas, err := ExtractVirtualInterfaces(page)
	return len(vas) == 0, err
}

// ExtractVirtualInterfaces interprets a VirtualInterfacePage as a slice of
// VirtualInterfaces.
func ExtractVirtualInterfaces(page pagination.Page) ([]VirtualInterface, error) {
	var response struct {
		VirtualInterfaces []VirtualInterface `mapstructure:"virtual_interfaces"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.VirtualInterfaces, nil
}

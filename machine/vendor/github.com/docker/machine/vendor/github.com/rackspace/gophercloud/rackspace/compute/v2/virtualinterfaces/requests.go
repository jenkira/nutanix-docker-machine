package virtualinterfaces

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of
// VirtualInterfaces associated with a given server.
func List(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	url := listURL(client, serverID)
	createPage := func(r pagination.PageResult) pagination.Page {
		return VirtualInterfacePage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPage)
}

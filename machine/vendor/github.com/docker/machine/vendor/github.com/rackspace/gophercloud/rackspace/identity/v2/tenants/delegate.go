package tenants

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/pagination"
)

// List enumerates the Tenants to which the current token has access.
func List(client *gophercloud.ServiceClient, opts *os.ListOpts) pagination.Pager {
	return os.List(client, opts)
}

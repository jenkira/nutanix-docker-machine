package extensions

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/extensions"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns a Pager which allows you to iterate over the full collection of extensions.
// It accepts a ListOpts struct, which allows you to provide additional options.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	return os.List(c)
}

// Get retrieves a particular extension based on its unique alias.
func Get(c *gophercloud.ServiceClient, alias string) os.GetResult {
	return os.Get(c, alias)
}

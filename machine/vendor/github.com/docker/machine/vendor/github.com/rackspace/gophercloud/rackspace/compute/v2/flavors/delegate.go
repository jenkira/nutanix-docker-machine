package flavors

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
)

// ListDetail returns a Pager that allows you to iterate over a full collection
// of flavors.
func ListDetail(client *gophercloud.ServiceClient, opts os.ListOptsBuilder) pagination.Pager {
	return os.ListDetail(client, opts)
}

// Get retrieves details of a single flavor. Use Extract to convert its result
// into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) os.GetResult {
	return os.Get(client, id)
}

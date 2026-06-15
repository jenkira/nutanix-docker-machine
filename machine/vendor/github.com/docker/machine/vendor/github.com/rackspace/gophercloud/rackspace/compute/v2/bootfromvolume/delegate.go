package bootfromvolume

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
)

// Create requests the creation of a server from a volume, using the options
// provided.
func Create(c *gophercloud.ServiceClient, opts os.CreateOptsBuilder) os.CreateResult {
	return os.Create(c, opts)
}

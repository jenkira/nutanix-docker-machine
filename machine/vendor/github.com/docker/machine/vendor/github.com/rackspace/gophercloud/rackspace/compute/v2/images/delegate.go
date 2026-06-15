package images

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/pagination"
)

// ListDetail enumerates the available images.
func ListDetail(client *gophercloud.ServiceClient, opts os.ListOptsBuilder) pagination.Pager {
	return os.ListDetail(client, opts)
}

// Get acquires additional detail about a specific image by ID.
// Use the Extract method to obtain a Image.
func Get(client *gophercloud.ServiceClient, id string) os.GetResult {
	return os.Get(client, id)
}

// Delete deletes the specified server image.
func Delete(client *gophercloud.ServiceClient, id string) os.DeleteResult {
	return os.Delete(client, id)
}

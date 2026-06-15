package volumeattach

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

// List returns a Pager that allows you to iterate over a collection of
// VolumeAttachments associated with a given server.
func List(client *gophercloud.ServiceClient, serverID string) os.ListResult {
	return os.List(client, serverID)
}

// Create requests the creation of a new VolumeAttachment on the server.
func Create(client *gophercloud.ServiceClient, serverID string, opts os.CreateOptsBuilder) os.CreateResult {
	return os.Create(client, serverID, opts)
}

// Get returns data about a specific VolumeAttachment.
func Get(client *gophercloud.ServiceClient, serverID, aID string) os.GetResult {
	return os.Get(client, serverID, aID)
}

// Delete requests the deletion of a VolumeAttachment.
func Delete(client *gophercloud.ServiceClient, serverID, aID string) os.DeleteResult {
	return os.Delete(client, serverID, aID)
}

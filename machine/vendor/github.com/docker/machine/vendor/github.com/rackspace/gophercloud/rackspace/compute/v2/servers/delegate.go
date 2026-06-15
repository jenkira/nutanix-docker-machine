package servers

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

// List makes a request against the API to list servers accessible to you.
func List(client *gophercloud.ServiceClient, opts os.ListOptsBuilder) pagination.Pager {
	return os.List(client, opts)
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(client *gophercloud.ServiceClient, opts os.CreateOptsBuilder) os.CreateResult {
	return os.Create(client, opts)
}

// Delete requests that a server previously provisioned be removed from your
// account.
func Delete(client *gophercloud.ServiceClient, id string) os.DeleteResult {
	return os.Delete(client, id)
}

// Get returns data about a specific server, by ID.
func Get(client *gophercloud.ServiceClient, id string) os.GetResult {
	return os.Get(client, id)
}

// Update requests that various attributes of the indicated server be changed.
func Update(client *gophercloud.ServiceClient, id string, opts os.UpdateOptsBuilder) os.UpdateResult {
	return os.Update(client, id, opts)
}

// ChangeAdminPassword alters the administrator or root password for a specified
// server.
func ChangeAdminPassword(client *gophercloud.ServiceClient, id, newPassword string) os.ActionResult {
	return os.ChangeAdminPassword(client, id, newPassword)
}

// Reboot requests that a given server reboot.
// Two types exist (soft and hard): Soft and Hard.
func Reboot(client *gophercloud.ServiceClient, id string, how os.RebootMethod) os.ActionResult {
	return os.Reboot(client, id, how)
}

// Rebuild will reprovision the server according to the configuration options
// provided in the RebuildOpts struct.
func Rebuild(client *gophercloud.ServiceClient, id string, opts os.RebuildOptsBuilder) os.RebuildResult {
	return os.Rebuild(client, id, opts)
}

// Resize instructs the provider to change the Flavor of the server.
func Resize(client *gophercloud.ServiceClient, id string, opts os.ResizeOptsBuilder) os.ActionResult {
	return os.Resize(client, id, opts)
}

// ConfirmResize confirms a previous resize operation on a server.
func ConfirmResize(client *gophercloud.ServiceClient, id string) os.ActionResult {
	return os.ConfirmResize(client, id)
}

// RevertResize cancels a previous resize operation on a server.
func RevertResize(client *gophercloud.ServiceClient, id string) os.ActionResult {
	return os.RevertResize(client, id)
}

package roles

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/roles"
	"github.com/rackspace/gophercloud/pagination"
)

// List enumerates the Roles to which the current token has access.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return os.List(client)
}

// AddUser grants a role to a user in a specific tenant.
func AddUser(client *gophercloud.ServiceClient, tenantID, userID, roleID string) os.AddUserRoleResult {
	return os.AddUser(client, tenantID, userID, roleID)
}

// DeleteUser deletes a role from a user in a specific tenant.
func DeleteUser(client *gophercloud.ServiceClient, tenantID, userID, roleID string) os.DeleteUserRoleResult {
	return os.DeleteUser(client, tenantID, userID, roleID)
}

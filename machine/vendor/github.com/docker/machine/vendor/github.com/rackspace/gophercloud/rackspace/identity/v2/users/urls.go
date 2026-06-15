package users

import "github.com/rackspace/gophercloud"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("users")
}

func listRolesURL(c *gophercloud.ServiceClient, tenantID, userID string) string {
	return c.ServiceURL("tenants", tenantID, "users", userID, "roles")
}

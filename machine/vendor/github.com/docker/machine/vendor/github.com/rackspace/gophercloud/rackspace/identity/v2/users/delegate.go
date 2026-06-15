package users

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List enumerates the Users.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, rootURL(client), func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.SinglePageBase(r)}
	})
}

// ListRoles returns a list of roles for a particular user and tenant.
func ListRoles(client *gophercloud.ServiceClient, tenantID, userID string) pagination.Pager {
	return pagination.NewPager(client, listRolesURL(client, tenantID, userID), func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.SinglePageBase(r)}
	})
}

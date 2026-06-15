package acl

import "github.com/rackspace/gophercloud"

func rootURL(c *gophercloud.ServiceClient, lbID int) string {
	return c.ServiceURL("loadbalancers", strconv.Itoa(lbID), "accesslist")
}

func resourceURL(c *gophercloud.ServiceClient, lbID, itemID int) string {
	return c.ServiceURL("loadbalancers", strconv.Itoa(lbID), "accesslist", strconv.Itoa(itemID))
}

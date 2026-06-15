package networks

import "github.com/rackspace/gophercloud"

func getListURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-networksv2")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("os-networksv2", id)
}

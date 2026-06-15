package lbs

import (
	"strconv"

	"github.com/rackspace/gophercloud"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("loadbalancers")
}

func resourceURL(c *gophercloud.ServiceClient, id int) string {
	return c.ServiceURL("loadbalancers", strconv.Itoa(id))
}

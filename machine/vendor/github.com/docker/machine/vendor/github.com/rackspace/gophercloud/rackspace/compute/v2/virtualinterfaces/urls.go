package virtualinterfaces

import "github.com/rackspace/gophercloud"

func listURL(c *gophercloud.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, "os-virtual-interfacesv2")
}

package networks

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of Networks.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewSinglePageBase(pagination.LastHTTPResponse(client.Get(getListURL(client), nil, nil)))
}

// Get returns data about a previously created Network.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = client.Get(getURL(client, id), &res.Body, nil)
	return res
}

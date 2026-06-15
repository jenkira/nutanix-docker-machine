package lbs

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/nodes"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/vips"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToLBListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	NodeAddress string `q:"nodeaddress"`
	ID          []int
}

// ToLBListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLBListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// load balancers.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)

	if opts != nil {
		query, err := opts.ToLBListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return LBPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a common interface.
type CreateOptsBuilder interface {
	ToLBCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Required - the name of the load balancer to create. The name must be 128
	// characters or less in length, and all UTF-8 characters are valid.
	Name string

	// Optional - the list of nodes to create with the LB.
	Nodes []nodes.Node

	// Required - the protocol of the service that is being load balanced.
	Protocol string

	// Optional - enables or disables Half-Closed support for the load balancer.
	HalfClosed gophercloud.EnabledState

	// Optional - specifies the algorithm that defines how traffic should be
	// directed between back-end nodes.
	Algorithm string

	// Optional - current connection logging configuration.
	ConnectionLogging *ConnectionLogging

	// Optional - the type of health monitor check to perform to ensure that the
	// service is performing properly.
	//HealthMonitor *HealthMonitor

	// Required for HTTPS protocol - configures a load balancer to terminate
	// HTTPS and begin HTTP at the front end.
	// VirtualIps []vips.VirtualIp
	VirtualIps []vips.VirtualIp

	// Optional - specifies a limit on the number of connections per IP address
	// to help mitigate malicious or errant clients.
	//AccessList *AccessList

	// Optional - the access list management feature allows fine-grained network
	// access controls to be applied to the load balancer virtual IP address.
	//RateLimit *RateLimit

	// Optional - the timeout value for the load balancer and communications
	// with its nodes.
	Timeout int

	// Optional - specifies whether or not the load balancer can target multiple
	// ports on a node.
	Port int
}

// ToLBCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToLBCreateMap() (map[string]interface{}, error) {
	lb := make(map[string]interface{})

	if opts.Name == "" {
		return lb, errors.New("Name is a required field")
	}
	if opts.Protocol == "" {
		return lb, errors.New("Protocol is a required field")
	}
	if len(opts.VirtualIps) == 0 {
		return lb, errors.New("VirtualIps is a required field")
	}

	lb["name"] = opts.Name
	lb["protocol"] = opts.Protocol

	if len(opts.Nodes) > 0 {
		nodeList := []map[string]interface{}{}
		for _, n := range opts.Nodes {
			nodeList = append(nodeList, map[string]interface{}{
				"address":   n.Address,
				"port":      n.Port,
				"condition": n.Condition,
			})
		}
		lb["nodes"] = nodeList
	}

	vipList := []map[string]interface{}{}
	for _, v := range opts.VirtualIps {
		vipList = append(vipList, map[string]interface{}{"type": v.Type})
	}
	lb["virtualIps"] = vipList

	if opts.HalfClosed != "" {
		lb["halfClosed"] = opts.HalfClosed
	}
	if opts.Algorithm != "" {
		lb["algorithm"] = opts.Algorithm
	}
	if opts.ConnectionLogging != nil {
		lb["connectionLogging"] = map[string]interface{}{"enabled": opts.ConnectionLogging.Enabled}
	}
	if opts.Timeout != 0 {
		lb["timeout"] = opts.Timeout
	}
	if opts.Port != 0 {
		lb["port"] = opts.Port
	}

	return map[string]interface{}{"loadBalancer": lb}, nil
}

// Create is the operation responsible for asynchronously provisioning a new
// load balancer based on the configuration defined in CreateOpts. Once the
// request is validated and progress has started on the provisioning process, a
// response object will be returned. The object will contain a unique ID and
// status of BUILD.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToLBCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(rootURL(c), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

// Get is the operation responsible for showing details of a load balancer.
func Get(c *gophercloud.ServiceClient, id int) GetResult {
	var res GetResult
	_, res.Err = c.Get(resourceURL(c, id), &res.Body, nil)
	return res
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package.
type UpdateOptsBuilder interface {
	ToLBUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Optional - the name of the load balancer to create. The name must be 128
	// characters or less in length, and all UTF-8 characters are valid.
	Name string

	// Optional - the protocol of the service that is being load balanced.
	Protocol string

	// Optional - enables or disables Half-Closed support for the load balancer.
	HalfClosed gophercloud.EnabledState

	// Optional - specifies the algorithm that defines how traffic should be
	// directed between back-end nodes.
	Algorithm string

	// Optional - current connection logging configuration.
	ConnectionLogging *ConnectionLogging

	// Optional - the timeout value for the load balancer and communications
	// with its nodes.
	Timeout int

	// Optional - specifies whether or not the load balancer can target multiple
	// ports on a node.
	Port int
}

// ToLBUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToLBUpdateMap() (map[string]interface{}, error) {
	lb := make(map[string]interface{})

	if opts.Name != "" {
		lb["name"] = opts.Name
	}
	if opts.Protocol != "" {
		lb["protocol"] = opts.Protocol
	}
	if opts.HalfClosed != "" {
		lb["halfClosed"] = opts.HalfClosed
	}
	if opts.Algorithm != "" {
		lb["algorithm"] = opts.Algorithm
	}
	if opts.ConnectionLogging != nil {
		lb["connectionLogging"] = map[string]interface{}{"enabled": opts.ConnectionLogging.Enabled}
	}
	if opts.Timeout != 0 {
		lb["timeout"] = opts.Timeout
	}
	if opts.Port != 0 {
		lb["port"] = opts.Port
	}

	return map[string]interface{}{"loadBalancer": lb}, nil
}

// Update is the operation responsible for asynchronously updating the
// attributes of a load balancer.
func Update(c *gophercloud.ServiceClient, id int, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToLBUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Put(resourceURL(c, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

// Delete is the operation responsible for permanently deleting a load balancer.
func Delete(c *gophercloud.ServiceClient, id int) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, id), &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

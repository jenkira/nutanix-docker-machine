package lbs

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/nodes"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/vips"
)

// Protocol represents the network protocol which the LB accepts.
type Protocol struct {
	Name string
	Port int
}

// ConnectionLogging - whether connection logging is enabled for the LB.
type ConnectionLogging struct {
	Enabled bool
}

// ContentCaching represents the content caching configuration of the LB.
type ContentCaching struct {
	Enabled bool
}

// SourceAddresses represents the source addresses of a load balancer.
type SourceAddresses struct {
	IPv4Servicenet string `mapstructure:"ipv4Servicenet"`
	IPv4Public     string `mapstructure:"ipv4Public"`
	IPv6Public     string `mapstructure:"ipv6Public"`
}

// Cluster represents the cluster a load balancer belongs to.
type Cluster struct {
	Name string `mapstructure:"name"`
}

// LoadBalancer represents a load balancer API resource.
type LoadBalancer struct {
	// The unique ID for the load balancer.
	ID int `mapstructure:"id"`

	// Human-readable name for the load balancer. Does not have to be unique.
	Name string `mapstructure:"name"`

	// Represents the port on which the client traffic is received.
	Port int `mapstructure:"port"`

	// Defines the protocol of the service that is being load balanced.
	Protocol string `mapstructure:"protocol"`

	// Specifies the algorithm that defines how traffic should be directed
	// between back-end nodes.
	Algorithm string `mapstructure:"algorithm"`

	// The current status of the load balancer.
	Status string `mapstructure:"status"`

	// The number of load balancer nodes.
	NodeCount int `mapstructure:"nodeCount"`

	// Nodes that are part of this load balancer.
	Nodes []nodes.Node `mapstructure:"nodes"`

	// Virtual IPs associated with this load balancer.
	VirtualIps []vips.VirtualIp `mapstructure:"virtualIps"`

	// The cluster the LB belongs to.
	Cluster Cluster `mapstructure:"cluster"`

	// Datetime when the LB was created.
	Created gophercloud.JSONRFC3339MilliNoZ `mapstructure:"created"`

	// Datetime when the LB was updated.
	Updated gophercloud.JSONRFC3339MilliNoZ `mapstructure:"updated"`

	// Datetime when the LB was updated.
	SourceAddresses SourceAddresses `mapstructure:"sourceAddresses"`

	HalfClosed       bool              `mapstructure:"halfClosed"`
	Timeout          int               `mapstructure:"timeout"`
	HttpsRedirect    bool              `mapstructure:"httpsRedirect"`
	ConnectionLogging ConnectionLogging `mapstructure:"connectionLogging"`
	ContentCaching   ContentCaching    `mapstructure:"contentCaching"`
}

// LBPage is the page returned by a pager when traversing over a collection
// of LBs.
type LBPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is used to step through a paginated list of LBs.
func (page LBPage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"loadBalancers_links"`
	}

	var r resp
	err := gophercloud.DecodeResponse(page, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// IsEmpty checks whether an LBPage struct is empty.
func (page LBPage) IsEmpty() (bool, error) {
	is, err := ExtractLBs(page)
	return len(is) == 0, err
}

// ExtractLBs accepts a Page struct, specifically an LBPage struct, and extracts
// the elements into a slice of LoadBalancer structs.
func ExtractLBs(page pagination.Page) ([]LoadBalancer, error) {
	var response struct {
		LoadBalancers []LoadBalancer `mapstructure:"loadBalancers"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancers, nil
}

// CommonResult represents the result of a Create or a Get operation.
type CommonResult struct {
	gophercloud.SingletonResult
}

// Extract interprets CommonResults as a LoadBalancer.
func (r CommonResult) Extract() (*LoadBalancer, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		LB *LoadBalancer `mapstructure:"loadBalancer"`
	}

	err := gophercloud.DecodeResponse(r.Body, &response)
	if err != nil {
		return nil, err
	}

	return response.LB, nil
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	CommonResult
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	CommonResult
}

// UpdateResult represents the result of an Update operation.
type UpdateResult struct {
	gophercloud.ErrResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

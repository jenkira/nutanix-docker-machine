package acl

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Type represents the access rule type.
type Type string

const (
	// ALLOW represents an allowlist entry.
	ALLOW Type = "ALLOW"

	// DENY represents a blocklist entry.
	DENY Type = "DENY"
)

// AccessListItem represents a single item in a load balancer's Access Control List.
type AccessListItem struct {
	ID      int    `mapstructure:"id"`
	Address string `mapstructure:"address"`
	Type    Type   `mapstructure:"type"`
}

// AccessListPage is the page returned by a pager when traversing over a collection
// of access control list items.
type AccessListPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether an AccessListPage struct is empty.
func (p AccessListPage) IsEmpty() (bool, error) {
	is, err := ExtractAccessList(p)
	return len(is) == 0, err
}

// ExtractAccessList accepts a Page struct, specifically an AccessListPage struct, and
// extracts the elements into a slice of AccessListItem structs.
func ExtractAccessList(page pagination.Page) ([]AccessListItem, error) {
	var response struct {
		Items []AccessListItem `mapstructure:"accessList"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

// CommonResult represents the result of a Create or Delete operation.
type CommonResult struct {
	gophercloud.ErrResult
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	CommonResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	CommonResult
}

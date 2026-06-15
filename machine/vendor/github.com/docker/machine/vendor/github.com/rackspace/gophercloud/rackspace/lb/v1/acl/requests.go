package acl

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List is the operation responsible for returning a paginated collection of
// network items that define access to the load balancer.
func List(c *gophercloud.ServiceClient, lbID int) pagination.Pager {
	url := rootURL(c, lbID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AccessListPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a common interface.
type CreateOptsBuilder interface {
	ToAccessListCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts []struct {
	// Required - the IP address for the ACL item.
	Address string

	// Required - the type of rule (ALLOW or DENY)
	Type Type
}

// ToAccessListCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToAccessListCreateMap() (map[string]interface{}, error) {
	var accessList []map[string]interface{}

	for _, o := range opts {
		if o.Address == "" {
			return nil, errors.New("Address is a required field")
		}
		if o.Type == "" {
			return nil, errors.New("Type is a required field")
		}
		acl := map[string]interface{}{
			"address": o.Address,
			"type":    o.Type,
		}
		accessList = append(accessList, acl)
	}

	return map[string]interface{}{"accessList": accessList}, nil
}

// Create is the operation responsible for creating or appending items to an
// access control list.
func Create(c *gophercloud.ServiceClient, lbID int, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToAccessListCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(rootURL(c, lbID), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

// Delete is the operation responsible for permanently deleting the entire
// contents of an access control list.
func Delete(c *gophercloud.ServiceClient, lbID int) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(rootURL(c, lbID), &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

// DeleteItem is the operation responsible for permanently deleting a single
// item in an access control list.
func DeleteItem(c *gophercloud.ServiceClient, lbID, itemID int) DeleteResult {
	var res DeleteResult
	_, res.Err = c.Delete(resourceURL(c, lbID, itemID), &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return res
}

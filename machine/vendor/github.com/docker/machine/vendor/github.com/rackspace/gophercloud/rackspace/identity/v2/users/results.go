package users

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// User represents a user resource.
type User struct {
	// UUID for the user
	ID string `mapstructure:"id"`

	// Human name for the user
	Name string `mapstructure:"name"`

	// The username for the user
	Username string `mapstructure:"username"`

	// The email address for the user
	Email string `mapstructure:"email"`

	// Whether or not the user is enabled
	Enabled bool `mapstructure:"enabled"`

	// The tenant ID for the user
	TenantID string `mapstructure:"tenant_id"`
}

// Role represents a role resource.
type Role struct {
	// UUID for the role
	ID string `mapstructure:"id"`

	// Human name for the role
	Name string `mapstructure:"name"`

	// Description of the role
	Description string `mapstructure:"description"`

	// The Tenant ID associated with this role
	TenantID string `mapstructure:"tenantId"`
}

// UserPage stores a page of Users.
type UserPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a UserPage is empty.
func (page UserPage) IsEmpty() (bool, error) {
	us, err := ExtractUsers(page)
	return len(us) == 0, err
}

// ExtractUsers interprets a UserPage as a slice of Users.
func ExtractUsers(page pagination.Page) ([]User, error) {
	var response struct {
		Users []User `mapstructure:"users"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.Users, nil
}

// RolePage stores a page of Roles.
type RolePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a RolePage is empty.
func (page RolePage) IsEmpty() (bool, error) {
	rs, err := ExtractRoles(page)
	return len(rs) == 0, err
}

// ExtractRoles interprets a RolePage as a slice of Roles.
func ExtractRoles(page pagination.Page) ([]Role, error) {
	var response struct {
		Roles []Role `mapstructure:"roles"`
	}

	err := gophercloud.DecodeResponse(page, &response)
	if err != nil {
		return nil, err
	}

	return response.Roles, nil
}

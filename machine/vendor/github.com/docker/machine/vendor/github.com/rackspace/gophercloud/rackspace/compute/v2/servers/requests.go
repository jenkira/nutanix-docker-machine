package servers

import (
	"errors"

	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// CreateOpts specifies all of the options that Rackspace accepts in its Create
// call, including those from the base OpenStack call.
type CreateOpts struct {
	// Name [required] is the name to assign to the newly launched server.
	Name string

	// ImageRef [required] is the ID or full URL to the image that contains the
	// server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageRef string

	// FlavorRef [required] is the ID or full URL to the flavor that describes the
	// server's specs.
	FlavorRef string

	// SecurityGroups [optional] lists the names of the security groups to which
	// this server should belong.
	SecurityGroups []string

	// UserData [optional] contains configuration information or scripts to use
	// upon launch. Create will base64-encode it for you.
	UserData []byte

	// AvailabilityZone [optional] is the name of the place where the server
	// will be created.
	AvailabilityZone string

	// Networks [optional] dictates how this server will be attached to available
	// networks.
	Networks []os.Network

	// Metadata [optional] contains key-value pairs (up to 255 bytes each)
	// to attach to the server.
	Metadata map[string]string

	// Personality [optional] includes the path and contents of a file to inject
	// into the server at launch.
	Personality os.FileContents

	// ConfigDrive [optional] enables metadata injection through a configuration
	// drive.
	ConfigDrive bool

	// AdminPass [optional] sets the root user password. If not set, a randomly
	// generated password will be created and returned in the response.
	AdminPass string

	// DiskConfig [optional] controls how the disk is partitioned when the server
	// is created.
	DiskConfig os.DiskConfig
}

// ToServerCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]interface{}, error) {
	server := make(map[string]interface{})

	if opts.Name == "" {
		return server, errors.New("Name is required")
	}
	server["name"] = opts.Name

	if opts.ImageRef == "" {
		return server, errors.New("ImageRef is required")
	}
	server["imageRef"] = opts.ImageRef

	if opts.FlavorRef == "" {
		return server, errors.New("FlavorRef is required")
	}
	server["flavorRef"] = opts.FlavorRef

	if opts.UserData != nil {
		server["user_data"] = opts.UserData
	}
	if opts.AvailabilityZone != "" {
		server["availability_zone"] = opts.AvailabilityZone
	}
	if len(opts.Networks) > 0 {
		var networks []map[string]interface{}
		for _, net := range opts.Networks {
			networks = append(networks, map[string]interface{}{"uuid": net.UUID, "port": net.Port})
		}
		server["networks"] = networks
	}
	if opts.Metadata != nil {
		server["metadata"] = opts.Metadata
	}
	if len(opts.Personality) > 0 {
		server["personality"] = opts.Personality
	}
	if opts.ConfigDrive {
		server["config_drive"] = opts.ConfigDrive
	}
	if opts.AdminPass != "" {
		server["adminPass"] = opts.AdminPass
	}
	if opts.DiskConfig != "" {
		server["OS-DCF:diskConfig"] = opts.DiskConfig
	}
	if len(opts.SecurityGroups) > 0 {
		securityGroups := make([]map[string]interface{}, len(opts.SecurityGroups))
		for i, sg := range opts.SecurityGroups {
			securityGroups[i] = map[string]interface{}{"name": sg}
		}
		server["security_groups"] = securityGroups
	}

	return map[string]interface{}{"server": server}, nil
}

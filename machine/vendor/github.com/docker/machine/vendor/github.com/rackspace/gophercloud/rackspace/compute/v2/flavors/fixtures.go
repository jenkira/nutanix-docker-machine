// +build fixtures

package flavors

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// ListOutput is a sample response to a ListDetail call.
const ListOutput = `
{
  "flavors": [
    {
      "OS-FLV-DISABLED:disabled": false,
      "disk": 1,
      "OS-FLV-EXT-DATA:ephemeral": 0,
      "os-flavor-access:is_public": true,
      "id": "1",
      "links": [
        {
          "href": "http://openstack.example.com/v2/openstack/flavors/1",
          "rel": "self"
        },
        {
          "href": "http://openstack.example.com/openstack/flavors/1",
          "rel": "bookmark"
        }
      ],
      "name": "m1.tiny",
      "ram": 512,
      "swap": "",
      "vcpus": 1
    },
    {
      "OS-FLV-DISABLED:disabled": false,
      "disk": 20,
      "OS-FLV-EXT-DATA:ephemeral": 0,
      "os-flavor-access:is_public": true,
      "id": "2",
      "links": [
        {
          "href": "http://openstack.example.com/v2/openstack/flavors/2",
          "rel": "self"
        },
        {
          "href": "http://openstack.example.com/openstack/flavors/2",
          "rel": "bookmark"
        }
      ],
      "name": "m1.small",
      "ram": 2048,
      "swap": "",
      "vcpus": 1
    },
    {
      "OS-FLV-DISABLED:disabled": false,
      "disk": 40,
      "OS-FLV-EXT-DATA:ephemeral": 0,
      "os-flavor-access:is_public": true,
      "id": "3",
      "links": [
        {
          "href": "http://openstack.example.com/v2/openstack/flavors/3",
          "rel": "self"
        },
        {
          "href": "http://openstack.example.com/openstack/flavors/3",
          "rel": "bookmark"
        }
      ],
      "name": "m1.medium",
      "ram": 4096,
      "swap": "",
      "vcpus": 2
    },
    {
      "OS-FLV-DISABLED:disabled": false,
      "disk": 80,
      "OS-FLV-EXT-DATA:ephemeral": 0,
      "os-flavor-access:is_public": true,
      "id": "4",
      "links": [
        {
          "href": "http://openstack.example.com/v2/openstack/flavors/4",
          "rel": "self"
        },
        {
          "href": "http://openstack.example.com/openstack/flavors/4",
          "rel": "bookmark"
        }
      ],
      "name": "m1.large",
      "ram": 8192,
      "swap": "",
      "vcpus": 4
    },
    {
      "OS-FLV-DISABLED:disabled": false,
      "disk": 160,
      "OS-FLV-EXT-DATA:ephemeral": 0,
      "os-flavor-access:is_public": true,
      "id": "5",
      "links": [
        {
          "href": "http://openstack.example.com/v2/openstack/flavors/5",
          "rel": "self"
        },
        {
          "href": "http://openstack.example.com/openstack/flavors/5",
          "rel": "bookmark"
        }
      ],
      "name": "m1.xlarge",
      "ram": 16384,
      "swap": "",
      "vcpus": 8
    }
  ]
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
  "flavor": {
    "OS-FLV-DISABLED:disabled": false,
    "disk": 1,
    "OS-FLV-EXT-DATA:ephemeral": 0,
    "os-flavor-access:is_public": true,
    "id": "1",
    "links": [
      {
        "href": "http://openstack.example.com/v2/openstack/flavors/1",
        "rel": "self"
      },
      {
        "href": "http://openstack.example.com/openstack/flavors/1",
        "rel": "bookmark"
      }
    ],
    "name": "m1.tiny",
    "ram": 512,
    "swap": "",
    "vcpus": 1
  }
}
`

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/flavors/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetSuccessfully configures the test server to respond to a Get request.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/flavors/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

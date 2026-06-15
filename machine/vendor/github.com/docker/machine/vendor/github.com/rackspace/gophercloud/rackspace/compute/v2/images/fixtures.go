// +build fixtures

package images

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
  "images": [
    {
      "status": "ACTIVE",
      "updated": "2014-09-23T12:54:56Z",
      "links": [
        {
          "href": "http://nova.example.com/v2/images/f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
          "rel": "self"
        },
        {
          "href": "http://nova.example.com/images/f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
          "rel": "bookmark"
        },
        {
          "href": "http://glance.example.com/images/f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
          "rel": "alternate",
          "type": "application/vnd.openstack.image"
        }
      ],
      "id": "f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7",
      "OS-EXT-IMG-SIZE:size": 476704768,
      "name": "F17-x86_64-cfntools",
      "created": "2014-09-23T12:54:52Z",
      "minDisk": 0,
      "progress": 100,
      "minRam": 0,
      "metadata": {}
    },
    {
      "status": "ACTIVE",
      "updated": "2014-09-23T12:51:03Z",
      "links": [
        {
          "href": "http://nova.example.com/v2/images/f90f6034-2570-4974-8351-6b49732ef2eb",
          "rel": "self"
        },
        {
          "href": "http://nova.example.com/images/f90f6034-2570-4974-8351-6b49732ef2eb",
          "rel": "bookmark"
        },
        {
          "href": "http://glance.example.com/images/f90f6034-2570-4974-8351-6b49732ef2eb",
          "rel": "alternate",
          "type": "application/vnd.openstack.image"
        }
      ],
      "id": "f90f6034-2570-4974-8351-6b49732ef2eb",
      "OS-EXT-IMG-SIZE:size": 13167616,
      "name": "cirros-0.3.2-x86_64-disk",
      "created": "2014-09-23T12:51:03Z",
      "minDisk": 0,
      "progress": 100,
      "minRam": 0,
      "metadata": {}
    }
  ]
}
`

// HandleListSuccessfully configures the test server to respond to a List
// request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetSuccessfully configures the test server to respond to a Get
// request.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images/f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete
// request.
func HandleDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/images/f3e4a95d-1f4f-4989-97ce-f3a1fb8c04d7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})
}

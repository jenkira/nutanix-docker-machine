// +build fixtures

package acl

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// HandleGetSuccessfully configures the test server to respond to a Get request for a particular ACL.
func HandleGetSuccessfully(t *testing.T, lbID int) {
	th.Mux.HandleFunc(fmt.Sprintf("/loadbalancers/%d/accesslist", lbID), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `
{
  "accessList": [
    {
      "address": "206.160.163.21",
      "id": 23,
      "type": "DENY"
    },
    {
      "address": "206.160.163.22",
      "id": 24,
      "type": "DENY"
    },
    {
      "address": "206.160.163.23",
      "id": 25,
      "type": "DENY"
    },
    {
      "address": "206.160.163.24",
      "id": 26,
      "type": "DENY"
    }
  ]
}
`)
		case "DELETE":
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.WriteHeader(http.StatusAccepted)
		}
	})

	th.Mux.HandleFunc(fmt.Sprintf("/loadbalancers/%d/accesslist/23", lbID), func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleCreateSuccessfully configures the test server to respond to a Create request for ACL items.
func HandleCreateSuccessfully(t *testing.T, lbID int) {
	th.Mux.HandleFunc(fmt.Sprintf("/loadbalancers/%d/accesslist", lbID), func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})
}

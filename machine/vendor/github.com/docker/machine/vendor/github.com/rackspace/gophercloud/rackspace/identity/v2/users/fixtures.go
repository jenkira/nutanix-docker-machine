// +build fixtures

package users

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// HandleListSuccessfully sets up the test server to respond to a user List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
  "users": [
    {
      "id": "u1000",
      "name": "John Smith",
      "username": "jsmith",
      "email": "jsmith@example.com",
      "enabled": true,
      "tenant_id": "12345"
    },
    {
      "id": "u1001",
      "name": "Jane Smith",
      "username": "janesmith",
      "email": "janesmith@example.com",
      "enabled": false,
      "tenant_id": "12345"
    }
  ]
}
		`)
	})
}

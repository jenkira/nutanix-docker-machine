// +build fixtures

package servers

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
	"servers": [
		{
			"status": "BUILD",
			"updated": "2015-03-17T16:21:29Z",
			"hostId": "",
			"addresses": {},
			"links": [
				{
					"href": "http://nova.example.com/v2/servers/asdfasdfasdf",
					"rel": "self"
				},
				{
					"href": "http://nova.example.com/servers/asdfasdfasdf",
					"rel": "bookmark"
				}
			],
			"key_name": null,
			"image": {
				"id": "e19a734c-c7e6-4a7f-921d-021840967571",
				"links": [
					{
						"href": "http://nova.example.com/images/e19a734c-c7e6-4a7f-921d-021840967571",
						"rel": "bookmark"
					}
				]
			},
			"flavor": {
				"id": "1",
				"links": [
					{
						"href": "http://nova.example.com/flavors/1",
						"rel": "bookmark"
					}
				]
			},
			"id": "asdfasdfasdf",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"user_id": "fake",
			"name": "test",
			"created": "2015-03-17T16:21:29Z",
			"tenant_id": "aaaabbbbcccc",
			"os-extended-volumes:volumes_attached": [],
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 0,
			"config_drive": "",
			"metadata": {}
		}
	]
}
`

// HandleListSuccessfully sets up the test server to respond to a server List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

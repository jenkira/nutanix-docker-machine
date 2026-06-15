// +build fixtures

package lbs

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

// SingleLBBody is the canned request/response body for a single LB.
const SingleLBBody = `
{
  "loadBalancer": {
    "name": "a-new-loadbalancer",
    "id": 144,
    "protocol": "HTTP",
    "halfClosed": false,
    "algorithm": "RANDOM",
    "status": "ACTIVE",
    "timeout": 30,
    "cluster": {
      "name": "ztm-n01.staging1.lbaas.rackspace.net"
    },
    "nodes": [
      {
        "address": "10.1.1.1",
        "id": 143,
        "port": 80,
        "status": "ONLINE",
        "condition": "ENABLED",
        "weight": 1
      }
    ],
    "virtualIps": [
      {
        "address": "206.10.10.210",
        "id": 39,
        "type": "PUBLIC",
        "ipVersion": "IPV4"
      }
    ],
    "created": {
      "time": "2010-11-30T03:23:42Z"
    },
    "updated": {
      "time": "2010-11-30T03:23:44Z"
    },
    "sourceAddresses": {
      "ipv6Public": "2001:4801:79f1:1::1/64",
      "ipv4Servicenet": "10.0.0.0",
      "ipv4Public": "10.12.99.28"
    },
    "httpsRedirect": false,
    "connectionLogging": {
      "enabled": false
    },
    "contentCaching": {
      "enabled": false
    }
  }
}
`

// HandleGetSuccessfully configures the test server to respond to a Get request for LB 144.
func HandleGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers/144", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, SingleLBBody)
	})
}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
  "loadBalancers": [
    {
      "name": "a-new-loadbalancer",
      "id": 144,
      "protocol": "HTTP",
      "halfClosed": false,
      "algorithm": "RANDOM",
      "status": "ACTIVE",
      "timeout": 30,
      "cluster": {
        "name": "ztm-n01.staging1.lbaas.rackspace.net"
      },
      "nodes": [
        {
          "address": "10.1.1.1",
          "id": 143,
          "port": 80,
          "status": "ONLINE",
          "condition": "ENABLED",
          "weight": 1
        }
      ],
      "virtualIps": [
        {
          "address": "206.10.10.210",
          "id": 39,
          "type": "PUBLIC",
          "ipVersion": "IPV4"
        }
      ],
      "created": {
        "time": "2010-11-30T03:23:42Z"
      },
      "updated": {
        "time": "2010-11-30T03:23:44Z"
      },
      "sourceAddresses": {
        "ipv6Public": "2001:4801:79f1:1::1/64",
        "ipv4Servicenet": "10.0.0.0",
        "ipv4Public": "10.12.99.28"
      },
      "httpsRedirect": false,
      "connectionLogging": {
        "enabled": false
      },
      "contentCaching": {
        "enabled": false
      }
    }
  ]
}
`)
	})
}

// HandleCreateSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, SingleLBBody)
	})
}

// HandleUpdateSuccessfully configures the test server to respond to an Update request for LB 144.
func HandleUpdateSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers/144", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete request for LB 144.
func HandleDeleteSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers/144", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}

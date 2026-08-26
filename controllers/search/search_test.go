package search

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

const testSearchSubnetsOnlyOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "subnets": {
      "code": 200,
      "data": [
        {
          "id": "3",
          "subnet": "168099840",
          "mask": "24",
		  "description": "Customer 1",
          "sectionId": "1",
          "masterSubnetId": "2",
          "allowRequests": "1",
          "vlanId": "0",
          "showName": "1",
          "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
          "pingSubnet": "0",
          "discoverSubnet": "0",
          "isFolder": "0",
          "isFull": "0"
        }
      ]
    }
  }
}
`

var testSearchSubnetsOnlyOutputExpected = []subnets.Subnet{
	{
		ID:             3,
		SubnetAddress:  "10.5.0.0",
		Mask:           24,
		SectionID:      1,
		Description:    "Customer 1",
		MasterSubnetID: 2,
		AllowRequests:  true,
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
}

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpOKTestServer(output string) *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, output, http.StatusOK)
	})
}

func fullSessionConfig() *session.Session {
	return &session.Session{
		Config: phpipam.Config{
			AppID:    "0123456789abcdefgh",
			Password: "changeit",
			Username: "nobody",
		},
		Token: session.Token{
			String: "foobarbazboop",
		},
	}
}

func TestSearch(t *testing.T) {
	ts := httpOKTestServer(testSearchSubnetsOnlyOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testSearchSubnetsOnlyOutputExpected
	actual, err := client.SearchSubnets("Customer 1")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

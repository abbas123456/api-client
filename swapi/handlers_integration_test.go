package swapi_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/kynrai/api-client/common/http"
	"github.com/kynrai/api-client/swapi"
)

func TestGetPersonIntegrationHandler(t *testing.T) {
	t.Parallel()
	person1data, err := ioutil.ReadFile("testdata/person.golden.json")
	if err != nil {
		t.Fatal("could not load file for Get Person test", err)
	}
	person1 := swapi.Person{}
	if err := json.Unmarshal(person1data, &person1); err != nil {
		t.Fatal("invalid JSON data for Get Person test", err)
	}
	for _, tc := range []struct {
		name string
		id   string
		url  string
		want *swapi.Person
		err  *h.HTTPError
	}{
		{
			name: "happy path",
			id:   "1",
			url:  "https://swapi.co/api/",
			want: &person1,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sw := swapi.New(
				swapi.URL(tc.url),
			)
			server := httptest.NewServer(swapi.GetPerson(sw))
			defer server.Close()
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			q.Add("id", tc.id)
			req.URL.RawQuery = q.Encode()

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(b))
		})
	}
}

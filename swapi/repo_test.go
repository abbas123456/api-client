package swapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	chttp "github.com/kynrai/api-client/common/http"
	"github.com/kynrai/api-client/swapi"
)

// We use _test package so that we can simulate how external packages would see the API of this package
// _test is recognised by GO nativly and can exist in the same package as its parent.

// Always use table driven tests even for 1 test. Makes it easy to add new tests in future.
func TestGetPerson(t *testing.T) {
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
		name      string
		id        string
		url       string
		transport http.RoundTripper
		expected  *swapi.Person
		err       error
	}{
		{
			name: "happy path",
			id:   "1",
			url:  "https://swapi.co/api/",
			transport: chttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(bytes.NewBuffer(person1data)),
				}, nil
			}),
			expected: &person1,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sw := swapi.New(
				swapi.Transport(tc.transport),
				swapi.URL(tc.url),
				swapi.Token("some token"),
			)
			// Should set some timeout on the context
			p, err := sw.GetPerson(context.Background(), tc.id)
			if tc.err != nil && !reflect.DeepEqual(err, tc.err) {
				t.Fatal(pretty.Compare(err.Error(), tc.err.Error()))
			}
			if tc.expected != nil && !reflect.DeepEqual(p, tc.expected) {
				t.Fatal(pretty.Compare(p, tc.expected))
			}
		})
	}
}

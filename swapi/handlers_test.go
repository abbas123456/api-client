package swapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	h "github.com/kynrai/api-client/common/http"
	"github.com/kynrai/api-client/swapi"
)

func TestGetPersonHandler(t *testing.T) {
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
		name   string
		id     string
		getter swapi.PersonGetter
		want   *swapi.Person
		err    *h.HTTPError
	}{
		{
			name: "happy path",
			id:   "1",
			getter: swapi.PersonGetterFunc(func(ctx context.Context, id string) (*swapi.Person, error) {
				return &person1, nil
			}),
			want: &person1,
		},
		{
			name: "missing id",
			getter: swapi.PersonGetterFunc(func(ctx context.Context, id string) (*swapi.Person, error) {
				return &person1, nil
			}),
			err: &h.HTTPError{Code: 400},
		},
		{
			name: "getter err",
			id:   "1",
			getter: swapi.PersonGetterFunc(func(ctx context.Context, id string) (*swapi.Person, error) {
				return nil, errors.New("Boom!")
			}),
			err: &h.HTTPError{Code: 500},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := swapi.GetPerson(tc.getter)
			req := httptest.NewRequest(http.MethodGet, "http://myapi.com?id="+tc.id, nil)
			rw := httptest.NewRecorder()
			handler.ServeHTTP(rw, req)
			if tc.err != nil && tc.err.Code != rw.Code {
				t.Fatalf("status codes dont match, Got: %d, Want: %d", rw.Code, tc.err.Code)
			}
			if tc.want != nil {
				var p swapi.Person
				if err := json.NewDecoder(rw.Body).Decode(&p); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(&p, tc.want) {
					t.Fatal(pretty.Compare(&p, tc.want))
				}
			}
		})
	}
}

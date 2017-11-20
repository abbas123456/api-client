package swapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

var _ Repo = (*swapi)(nil)

type Repo interface {
	PersonGetter
}

type (
	PersonGetter interface {
		GetPerson(ctx context.Context, id string) (*Person, error)
	}
	PersonGetterFunc func(ctx context.Context, id string) (*Person, error)
)

func (f PersonGetterFunc) Get(ctx context.Context, id string) (*Person, error) {
	return f(ctx, id)
}

type swapi struct {
	url        string       // The base URL of the API
	token      string       // Used for API tokens
	httpClient *http.Client // Used for mocking API calls
}

type swapiOption func(*swapi)

func URL(url string) swapiOption {
	return func(s *swapi) {
		s.url = url
	}
}

func Token(token string) swapiOption {
	return func(s *swapi) {
		s.token = token
	}
}

func Transport(rt http.RoundTripper) swapiOption {
	return func(s *swapi) {
		s.httpClient.Transport = rt
	}
}

func New(opts ...swapiOption) *swapi {
	s := new(swapi)
	s.httpClient = &http.Client{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *swapi) GetPerson(ctx context.Context, id string) (*Person, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", s.url, "people", id), nil)
	if err != nil {
		return nil, err
	}
	// This where you set headers etc
	// req.Header.Set("Authorization", s.token)

	// This is how we used to set query params
	// q := req.URL.Query()
	// q.Add("key","someKey")
	// req.URL.RawQuery = q.Encode()

	// Set some timeout in the context
	res, err := s.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	var p Person
	// This streams the response body into p, if we use ioutil.ReadAll then it
	// allocates memory to read in the entire response before copying
	// into a struct with json.Unmarshal
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

package http

import "net/http"

// RoundTripperFunc is a RoundTripper transport used in a http.Client
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements RoundTripper for the RoundTripperFunc
func (r RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

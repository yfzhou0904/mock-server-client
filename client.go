package mockserver

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type MockClient struct {
	restyClient *resty.Client
}

// NewClient creates a new client provided its host and port
func NewClient(host string, port int) MockClient {
	return NewClientURL(fmt.Sprintf("http://%s:%d", host, port))
}

// NewClientURL creates a new client provided its URL
func NewClientURL(url string) MockClient {
	return MockClient{
		restyClient: resty.New().
			SetHostURL(url),
	}
}

// SetDebug enables or disables the debug
func (c MockClient) SetDebug(d bool) MockClient {
	c.restyClient.SetDebug(d)
	return c
}

// Verify checks if the mock server received requests matching the matcher.
func (c MockClient) Verify(matcher RequestMatcher, times Times) error {
	payload := map[string]interface{}{
		"httpRequest": matcher,
		"times":       times,
	}

	resp, err := c.restyClient.NewRequest().
		SetDoNotParseResponse(true).
		SetBody(payload).
		Put("/mockserver/verify")

	if err != nil {
		return errors.Wrap(err, "error calling verify endpoint")
	}
	if resp.StatusCode() == http.StatusAccepted {
		return nil
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.RawBody())

	return errors.Wrap(errors.New(buf.String()), "verification failed")
}

// Clear erases from the mock server all the requests matching the matcher.
func (c MockClient) Clear(matcher RequestMatcher) error {
	resp, err := c.restyClient.NewRequest().
		SetBody(matcher).
		Put("mockserver/clear?type=LOG")

	if err != nil {
		return errors.Wrap(err, "error calling clear endpoint (type=LOGS)")
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.Wrap(errors.New("status was expected to be 200"), "log clearing failed")
	}
	return nil
}

// VerifyAndClear checks if the mock server received requests matching the matcher
// and then erases from the logs the requests matching the matcher.
func (c MockClient) VerifyAndClear(matcher RequestMatcher, times Times) error {
	err_verify := c.Verify(matcher, times)
	err_clear := c.Clear(matcher)
	if err_verify != nil {
		return errors.Wrap(err_verify, "could not verify")
	}
	if err_clear != nil {
		return errors.Wrap(err_clear, "could not clear")
	}
	return nil
}

// VerifyAndClearByHeader checks if the mock server received requests matching the matcher
// and having the specified header name and value.
// It then erases from the logs the requests matching the same header name and value.
func (c MockClient) VerifyAndClearByHeader(headerName, headerValue string, matcher RequestMatcher, times Times) error {
	err_verify := c.Verify(matcher.WithHeader(headerName, headerValue), times)
	err_clear := c.Clear(RequestMatcher{}.WithHeader(headerName, headerValue))
	if err_verify != nil {
		return errors.Wrap(err_verify, "could not verify")
	}
	if err_clear != nil {
		return errors.Wrap(err_clear, "could not clear")
	}
	return nil
}

// Set a new Expectation in mock server with request and response
func (c MockClient) RegisterExpectation(expectation Expectation) error {
	_, err := c.restyClient.NewRequest().
		SetDoNotParseResponse(true).
		SetBody(expectation).
		Put("/mockserver/expectation")

	if err != nil {
		return errors.Wrap(err, "error calling SetExpectation endpoint")
	}

	return nil
}

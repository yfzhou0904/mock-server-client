package mockserver

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/palantir/stacktrace"
)

type MockClient struct {
	restyClient *resty.Client
}

func NewClient(host string, port int) MockClient {
	return MockClient{
		restyClient: resty.New().
			SetHostURL(fmt.Sprintf("http://%s:%d", host, port)),
	}
}

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
		return stacktrace.Propagate(err, "Error calling verify endpoint")
	}
	if resp.StatusCode() == http.StatusAccepted {
		return nil
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.RawBody())

	return stacktrace.Propagate(errors.New(buf.String()), "Verification failed")
}

func (c MockClient) Clear(matcher RequestMatcher) error {
	resp, err := c.restyClient.NewRequest().
		SetBody(matcher).
		Put("mockserver/clear?type=LOG")

	if err != nil {
		return stacktrace.Propagate(err, "Error calling clear endpoint (type=LOGS)")
	}

	if resp.StatusCode() != http.StatusOK {
		return stacktrace.Propagate(errors.New("Status was expected to be 200"), "Log clearing failed")
	}
	return nil
}

func (c MockClient) VerifyAndClear(matcher RequestMatcher, times Times) error {
	err_verify := c.Verify(matcher, times)
	err_clear := c.Clear(matcher)
	if err_verify != nil {
		return stacktrace.Propagate(err_verify, "Could not verify")
	}
	if err_clear != nil {
		return stacktrace.Propagate(err_clear, "Could not clear")
	}
	return nil
}

func (c MockClient) VerifyAndClearByHeader(headerName, headerValue string, matcher RequestMatcher, times Times) error {
	err_verify := c.Verify(matcher.WithHeader(headerName, headerValue), times)
	err_clear := c.Clear(RequestMatcher{}.WithHeader(headerName, headerValue))
	if err_verify != nil {
		return stacktrace.Propagate(err_verify, "Could not verify")
	}
	if err_clear != nil {
		return stacktrace.Propagate(err_clear, "Could not clear")
	}
	return nil
}

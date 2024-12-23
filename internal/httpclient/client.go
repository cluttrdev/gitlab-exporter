package httpclient

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	*retryablehttp.Client
}

func New() *Client {
	return &Client{
		Client: &retryablehttp.Client{
			Backoff:      retryHTTPBackoff,
			CheckRetry:   retryHTTPCheck,
			ErrorHandler: retryablehttp.PassthroughErrorHandler,
			HTTPClient:   cleanhttp.DefaultPooledClient(),
			RetryWaitMin: 100 * time.Millisecond,
			RetryWaitMax: 400 * time.Millisecond,
			RetryMax:     5,
		},
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Convert the request to be retryable.
	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	// Execute the request.
	resp, err := c.Client.Do(retryableReq)
	// If we got an error returned by standard library's `Do` method, unwrap it
	// otherwise we will wind up erroneously re-nesting the error.
	if _, ok := err.(*url.Error); ok {
		return resp, errors.Unwrap(err)
	}

	return resp, err
}

const (
	headerRateLimit = "RateLimit-Limit"
	headerRateReset = "RateLimit-Reset"
)

// retryHTTPCheck provides a callback for Client.CheckRetry which
// will retry both rate limit (429) and server (>= 500) errors.
func retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}

// retryHTTPBackoff provides a generic callback for Client.Backoff which
// will pass through all calls based on the status code of the response.
func retryHTTPBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	// Use the rate limit backoff function when we are rate limited.
	if resp != nil && resp.StatusCode == 429 {
		return rateLimitBackoff(min, max, attemptNum, resp)
	}

	// Set custom duration's when we experience a service interruption.
	min = 700 * time.Millisecond
	max = 900 * time.Millisecond

	return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
}

// rateLimitBackoff provides a callback for Client.Backoff which will use the
// RateLimit-Reset header to determine the time to wait. We add some jitter
// to prevent a thundering herd.
//
// min and max are mainly used for bounding the jitter that will be added to
// the reset time retrieved from the headers. But if the final wait time is
// less then min, min will be used instead.
func rateLimitBackoff(min, max time.Duration, _ int, resp *http.Response) time.Duration {
	// rnd is used to generate pseudo-random numbers.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// First create some jitter bounded by the min and max durations.
	jitter := time.Duration(rnd.Float64() * float64(max-min))

	if resp != nil {
		if v := resp.Header.Get(headerRateReset); v != "" {
			if reset, _ := strconv.ParseInt(v, 10, 64); reset > 0 {
				// Only update min if the given time to wait is longer.
				if wait := time.Until(time.Unix(reset, 0)); wait > min {
					min = wait
				}
			}
		}
	}

	return min + jitter
}

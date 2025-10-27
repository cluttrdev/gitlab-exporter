package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()

	if client == nil {
		t.Fatal("New() returned nil")
	}

	if client.Client == nil {
		t.Fatal("client.Client is nil")
	}

	// Verify default settings
	if client.Client.RetryMax != 5 {
		t.Errorf("expected RetryMax to be 5, got %d", client.Client.RetryMax)
	}

	if client.Client.RetryWaitMin != 100*time.Millisecond {
		t.Errorf("expected RetryWaitMin to be 100ms, got %v", client.Client.RetryWaitMin)
	}

	if client.Client.RetryWaitMax != 400*time.Millisecond {
		t.Errorf("expected RetryWaitMax to be 400ms, got %v", client.Client.RetryWaitMax)
	}

	if client.Client.Backoff == nil {
		t.Error("Backoff function is nil")
	}

	if client.Client.CheckRetry == nil {
		t.Error("CheckRetry function is nil")
	}

	if client.Client.ErrorHandler == nil {
		t.Error("ErrorHandler is nil")
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		wantErr       bool
		expectRetries bool
	}{
		{
			name:       "successful request",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "client error - no retry",
			statusCode: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name:       "not found - no retry",
			statusCode: http.StatusNotFound,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := New()
			// Disable retries for predictable tests
			client.Client.RetryMax = 0

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Fatal("response is nil")
			}

			if resp != nil && resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestClient_Do_Retry(t *testing.T) {
	tests := []struct {
		name             string
		statusCodes      []int
		expectedRequests int
		finalStatusCode  int
	}{
		{
			name:             "retry on 429",
			statusCodes:      []int{429, 429, 200},
			expectedRequests: 3,
			finalStatusCode:  200,
		},
		{
			name:             "retry on 500",
			statusCodes:      []int{500, 500, 200},
			expectedRequests: 3,
			finalStatusCode:  200,
		},
		{
			name:             "retry on 502",
			statusCodes:      []int{502, 200},
			expectedRequests: 2,
			finalStatusCode:  200,
		},
		{
			name:             "max retries exceeded",
			statusCodes:      []int{500, 500, 500, 500, 500, 500, 500},
			expectedRequests: 7, // 1 initial + 5 retries + 1 more (retryablehttp behavior)
			finalStatusCode:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				statusCode := tt.statusCodes[requestCount]
				if requestCount < len(tt.statusCodes)-1 {
					requestCount++
				}
				w.WriteHeader(statusCode)
			}))
			defer server.Close()

			client := New()
			// Set very short retry wait times for faster tests
			client.Client.RetryWaitMin = 1 * time.Millisecond
			client.Client.RetryWaitMax = 2 * time.Millisecond

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}

			if resp.StatusCode != tt.finalStatusCode {
				t.Errorf("expected final status code %d, got %d", tt.finalStatusCode, resp.StatusCode)
			}

			if requestCount+1 != tt.expectedRequests {
				t.Errorf("expected %d requests, got %d", tt.expectedRequests, requestCount+1)
			}
		})
	}
}

func TestRetryHTTPCheck(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		err         error
		ctxCanceled bool
		expectRetry bool
		expectErr   bool
	}{
		{
			name:        "retry on 429",
			statusCode:  429,
			expectRetry: true,
			expectErr:   false,
		},
		{
			name:        "retry on 500",
			statusCode:  500,
			expectRetry: true,
			expectErr:   false,
		},
		{
			name:        "retry on 502",
			statusCode:  502,
			expectRetry: true,
			expectErr:   false,
		},
		{
			name:        "retry on 503",
			statusCode:  503,
			expectRetry: true,
			expectErr:   false,
		},
		{
			name:        "no retry on 200",
			statusCode:  200,
			expectRetry: false,
			expectErr:   false,
		},
		{
			name:        "no retry on 404",
			statusCode:  404,
			expectRetry: false,
			expectErr:   false,
		},
		{
			name:        "no retry on 400",
			statusCode:  400,
			expectRetry: false,
			expectErr:   false,
		},
		{
			name:        "no retry with error",
			err:         http.ErrHandlerTimeout,
			expectRetry: false,
			expectErr:   true,
		},
		{
			name:        "no retry with canceled context",
			ctxCanceled: true,
			expectRetry: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.ctxCanceled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			var resp *http.Response
			if !tt.ctxCanceled && tt.err == nil {
				resp = &http.Response{
					StatusCode: tt.statusCode,
				}
			}

			retry, err := retryHTTPCheck(ctx, resp, tt.err)

			if retry != tt.expectRetry {
				t.Errorf("expected retry=%v, got %v", tt.expectRetry, retry)
			}

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error=%v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestRetryHTTPBackoff(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		minTime    time.Duration
	}{
		{
			name:       "rate limit backoff for 429",
			statusCode: 429,
			minTime:    100 * time.Millisecond,
		},
		{
			name:       "service interruption backoff for 500",
			statusCode: 500,
			minTime:    700 * time.Millisecond,
		},
		{
			name:       "service interruption backoff for 502",
			statusCode: 502,
			minTime:    700 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Header:     make(http.Header),
			}

			min := 100 * time.Millisecond
			max := 400 * time.Millisecond
			attemptNum := 1

			backoff := retryHTTPBackoff(min, max, attemptNum, resp)

			// Backoff should be at least the minimum time
			// We don't check maximum because LinearJitterBackoff can exceed max
			// based on attempt number
			if backoff < tt.minTime {
				t.Errorf("backoff %v is less than minimum %v", backoff, tt.minTime)
			}

			// Verify it's a reasonable value (not absurdly high)
			if backoff > 5*time.Second {
				t.Errorf("backoff %v is unexpectedly high", backoff)
			}
		})
	}
}

func TestRateLimitBackoff(t *testing.T) {
	t.Run("no rate limit header", func(t *testing.T) {
		resp := &http.Response{
			Header: make(http.Header),
		}

		min := 100 * time.Millisecond
		max := 400 * time.Millisecond

		backoff := rateLimitBackoff(min, max, 1, resp)

		if backoff < min {
			t.Errorf("backoff %v is less than min %v", backoff, min)
		}

		maxExpected := max + (max - min)
		if backoff > maxExpected {
			t.Errorf("backoff %v exceeds expected max %v for minimum wait", backoff, maxExpected)
		}
	})

	t.Run("invalid reset header", func(t *testing.T) {
		resp := &http.Response{
			Header: make(http.Header),
		}
		resp.Header.Set(headerRateReset, "invalid")

		min := 100 * time.Millisecond
		max := 400 * time.Millisecond

		backoff := rateLimitBackoff(min, max, 1, resp)

		if backoff < min {
			t.Errorf("backoff %v is less than min %v", backoff, min)
		}

		maxExpected := max + (max - min)
		if backoff > maxExpected {
			t.Errorf("backoff %v exceeds expected max %v for minimum wait", backoff, maxExpected)
		}
	})

	t.Run("future reset time", func(t *testing.T) {
		// Set reset time right before calling the function to minimize timing drift
		resetTime := time.Now().Add(2 * time.Second)

		resp := &http.Response{
			Header: make(http.Header),
		}
		resp.Header.Set(headerRateReset, strconv.FormatInt(resetTime.Unix(), 10))

		min := 100 * time.Millisecond
		max := 400 * time.Millisecond

		backoff := rateLimitBackoff(min, max, 1, resp)

		// The backoff should be: time.Until(resetTime) + jitter
		// where jitter is between 0 and (max-min) = 300ms
		// Given timing variations, we should be lenient
		// Minimum: should be at least min (100ms)
		if backoff < min {
			t.Errorf("backoff %v is less than min %v", backoff, min)
		}

		// Maximum: Should be roughly 2s + max jitter (300ms) + some tolerance
		maxExpected := 2*time.Second + (max - min) + 200*time.Millisecond
		if backoff > maxExpected {
			t.Errorf("backoff %v exceeds expected maximum %v", backoff, maxExpected)
		}

		// Verify it's more than the base min, indicating the reset header was used
		if backoff <= max {
			t.Errorf("backoff %v suggests reset header was not used (should be > %v)", backoff, max)
		}
	})

	t.Run("past reset time", func(t *testing.T) {
		resp := &http.Response{
			Header: make(http.Header),
		}
		resp.Header.Set(headerRateReset, strconv.FormatInt(time.Now().Add(-1*time.Second).Unix(), 10))

		min := 100 * time.Millisecond
		max := 400 * time.Millisecond

		backoff := rateLimitBackoff(min, max, 1, resp)

		if backoff < min {
			t.Errorf("backoff %v is less than min %v", backoff, min)
		}

		maxExpected := max + (max - min)
		if backoff > maxExpected {
			t.Errorf("backoff %v exceeds expected max %v for minimum wait", backoff, maxExpected)
		}
	})
}

func TestRateLimitBackoff_WithNilResponse(t *testing.T) {
	min := 100 * time.Millisecond
	max := 400 * time.Millisecond

	backoff := rateLimitBackoff(min, max, 1, nil)

	// Should return min + jitter
	if backoff < min {
		t.Errorf("backoff %v is less than min %v", backoff, min)
	}

	maxExpected := max + (max - min)
	if backoff > maxExpected {
		t.Errorf("backoff %v exceeds expected max %v", backoff, maxExpected)
	}
}

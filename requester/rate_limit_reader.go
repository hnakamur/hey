package requester

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

// RateLimitReader is a rate limitted reader.
type RateLimitReader struct {
	r       io.Reader
	limiter *rate.Limiter
}

var _ = io.Reader((*RateLimitReader)(nil))

// NewReader creates a rate limitted reader.
func NewRateLimitReader(r io.Reader, limiter *rate.Limiter) *RateLimitReader {
	return &RateLimitReader{r: r, limiter: limiter}
}

// Read implements io.Reader.
func (r *RateLimitReader) Read(p []byte) (n int, err error) {
	// Adapted from https://daichi.dev/posts/2023-01-19-golang-x-time-rate
	for n < len(p) {
		size := len(p[n:])
		if burst := r.limiter.Burst(); size > burst {
			size = burst
		}

		rn, err := r.r.Read(p[n : n+size])
		n += rn
		if err != nil {
			return n, err
		}

		if err := r.limiter.WaitN(context.Background(), rn); err != nil {
			return n, err
		}
	}
	return n, nil
}

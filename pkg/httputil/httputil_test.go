package httputil

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		url      string
		wantCode int
		wantErr  bool
	}{
		{
			name:     "success",
			url:      "https://www.google.com",
			wantCode: 200,
			wantErr:  false,
		},
		{
			name:     "bad request",
			url:      "https://www.google.com/not-found",
			wantCode: 404,
			wantErr:  false,
		},
		{
			name:     "host not found",
			url:      "https://not-found.google.com",
			wantCode: 0,
			wantErr:  true,
		},
		{
			name:     "no url",
			url:      "",
			wantCode: 0,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			resp, err := Get(ctx, tc.url)
			defer func() {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
			}()

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if resp.StatusCode != tc.wantCode {
				t.Errorf("unexpected status code: %d", resp.StatusCode)
			}
		})
	}
}

package crawler

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestCrawler_AcknowledgeRobotsTxt(t *testing.T) {
	type fields struct {
		client    *http.Client
		userAgent string
	}
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "robots.txt disallows crawling",
			fields: fields{
				client: NewTestClient(func(r *http.Request) *http.Response {
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("User-agent: *\nDisallow: /")),
						Header:     make(http.Header),
					}
					return resp
				}),
				userAgent: "fediverse-blahaj/0.0.1",
			},
			args: args{
				ctx: context.Background(),
				url: "https://misskey.takanotume24.com",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "robots.txt allows crawling",
			fields: fields{
				client: NewTestClient(func(r *http.Request) *http.Response {
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("User-agent: *\nDisallow: /api")),
						Header:     make(http.Header),
					}
					return resp
				}),
				userAgent: "fediverse-blahaj/0.0.1",
			},
			args: args{
				ctx: context.Background(),
				url: "https://misskey.takanotume24.com",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				client:    &retryablehttp.Client{HTTPClient: tt.fields.client},
				userAgent: tt.fields.userAgent,
			}
			got, err := c.acknowledgeRobotsTxt(tt.args.ctx, tt.args.url)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

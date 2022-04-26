package fanartapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is a fanart client.
type Client struct {
	ApiKey    string
	ClientKey string
	Transport http.RoundTripper
}

// New creates a new fanart client.
func New(opts ...Option) *Client {
	cl := &Client{}
	for _, o := range opts {
		o(cl)
	}
	return cl
}

// Do executes a request of typ, decoding results to v.
func (cl *Client) Do(ctx context.Context, typ string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://webservice.fanart.tv/v3/"+typ, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	if cl.ApiKey != "" {
		req.Header.Add("api-key", cl.ApiKey)
	}
	if cl.ClientKey != "" {
		req.Header.Add("client-key", cl.ClientKey)
	}
	httpClient := &http.Client{
		Transport: cl.Transport,
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d != 200", res.StatusCode)
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// Images retrieves the fanart images for type, id.
func (cl *Client) Images(ctx context.Context, typ Type, id string) (*ImagesResult, error) {
	return Images(typ, id).Do(ctx, cl)
}

// Latest retrieves the fanart latest for type, id.
func (cl *Client) Latest(ctx context.Context, typ Type) ([]LatestResult, error) {
	return Latest(typ).Do(ctx, cl)
}

// Option is a fanart client option.
type Option func(*Client)

// WithApiKey is a fanart client option to set the api key used.
func WithApiKey(apiKey string) Option {
	return func(cl *Client) {
		cl.ApiKey = apiKey
	}
}

// WithClientKey is a fanart client option to set the client key used.
func WithClientKey(clientKey string) Option {
	return func(cl *Client) {
		cl.ClientKey = clientKey
	}
}

// WithTransport is a fanart client option to set the http transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.Transport = transport
	}
}

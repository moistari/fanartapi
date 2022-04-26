package fanartapi

import (
	"context"
	"os"
	"testing"
)

func TestImages(t *testing.T) {
	tests := []struct {
		typ Type
		id  string
	}{
		{Movie, "tt0137523"},
		{Movie, "10195"},
		{Series, "75682"},
		{MusicArtist, "f4a31f0a-51dd-4fa7-986d-3095c40c5ed9"},
		{MusicAlbum, "9ba659df-5814-32f6-b95f-02b738698e7c"},
		{MusicLabel, "e832b688-546b-45e3-83e5-9f8db5dcde1d"},
	}
	cl := buildClient()
	for i, test := range tests {
		res, err := cl.Images(context.Background(), test.typ, test.id)
		if err != nil {
			t.Fatalf("test %d expected no error, got: %v", i, err)
		}
		t.Logf("%d: %q (%s)", i, res.Name, res.ID())
	}
}

func TestLatest(t *testing.T) {
	cl := buildClient()
	for i, typ := range []Type{Movie, Series, MusicArtist} {
		latest, err := cl.Latest(context.Background(), typ)
		if err != nil {
			t.Fatalf("test %d (%s) expected no error, got: %v", i, typ, err)
		}
		for _, res := range latest {
			t.Logf("%d: %s %q (%s)", i, typ, res.Name, res.ID())
		}
	}
}

func buildClient() *Client {
	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		apiKey = "6fa42b0ef3b5f3aab6a7edaa78675ac2"
	}
	opts := []Option{WithApiKey(apiKey)}
	if clientKey := os.Getenv("CLIENTKEY"); clientKey != "" {
		opts = append(opts, WithClientKey(clientKey))
	}
	return New(opts...)
}

package fanartapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Type is a query type.
type Type int

// Type values.
const (
	Movie Type = iota
	Series
	MusicArtist
	MusicAlbum
	MusicLabel
)

// String satisfies the fmt.Stringer interface.
func (typ Type) String() string {
	switch typ {
	case Movie:
		return "movie"
	case Series:
		return "series"
	case MusicArtist:
		return "artist"
	case MusicAlbum:
		return "album"
	case MusicLabel:
		return "label"
	}
	return fmt.Sprintf("Type(%d)", typ)
}

// ApiType returns the API type of the query type.
func (typ Type) ApiType() string {
	switch typ {
	case Movie:
		return "movies"
	case Series:
		return "tv"
	case MusicArtist:
		return "music"
	case MusicAlbum:
		return "music/albums"
	case MusicLabel:
		return "music/labels"
	}
	return "unknown"
}

// ImagesRequest is a fanart images request.
type ImagesRequest struct {
	Type  Type
	Query string
}

// Images creates a new images request.
func Images(typ Type, query string) *ImagesRequest {
	return &ImagesRequest{
		Type:  typ,
		Query: query,
	}
}

// Do executes the images request against the client.
func (req *ImagesRequest) Do(ctx context.Context, cl *Client) (*ImagesResult, error) {
	res := new(ImagesResult)
	if err := cl.Do(ctx, req.Type.ApiType()+"/"+req.Query, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ImagesResult is a fanart api images result.
type ImagesResult struct {
	Name             string                        `json:"name,omitempty"`
	FanartID         string                        `json:"id,omitempty"`
	ImdbID           string                        `json:"imdb_id,omitempty"`
	TmdbID           string                        `json:"tmdb_id,omitempty"`
	TvdbID           string                        `json:"thetvdb_id,omitempty"`
	MbID             string                        `json:"mbid_id,omitempty"`
	ClearArt         []Image                       `json:"clearart,omitempty"`
	ClearLogo        []Image                       `json:"clearlogo,omitempty"`
	HdClearArt       []Image                       `json:"hdclearart,omitempty"`
	HdLogo           []Image                       `json:"hdlogo,omitempty"`
	HdMovieClearArt  []Image                       `json:"hdmovieclearart,omitempty"`
	HdMovieLogo      []Image                       `json:"hdmovielogo,omitempty"`
	HdMusicLogo      []Image                       `json:"hdmusiclogo,omitempty"`
	HdtvLogo         []Image                       `json:"hdtvlogo,omitempty"`
	MovieArt         []Image                       `json:"movieart,omitempty"`
	CharacterArt     []Image                       `json:"characterart,omitempty"`
	MovieBackground  []Image                       `json:"moviebackground,omitempty"`
	MovieBanner      []Image                       `json:"moviebanner,omitempty"`
	MovieDisc        []Image                       `json:"moviedisc,omitempty"`
	MovieLogo        []Image                       `json:"movielogo,omitempty"`
	MoviePoster      []Image                       `json:"movieposter,omitempty"`
	MovieThumb       []Image                       `json:"moviethumb,omitempty"`
	ArtistThumb      []Image                       `json:"artistthumb,omitempty"`
	SeasonBanner     []Image                       `json:"seasonbanner,omitempty"`
	SeasonPoster     []Image                       `json:"seasonposter,omitempty"`
	SeasonThumb      []Image                       `json:"seasonthumb,omitempty"`
	ShowBackground   []Image                       `json:"showbackground,omitempty"`
	ArtistBackground []Image                       `json:"artistbackground,omitempty"`
	TvBanner         []Image                       `json:"tvbanner,omitempty"`
	TvPoster         []Image                       `json:"tvposter,omitempty"`
	TvThumb          []Image                       `json:"tvthumb,omitempty"`
	MusicLogo        []Image                       `json:"musiclogo,omitempty"`
	MusicBanner      []Image                       `json:"musicbanner,omitempty"`
	MusicLabel       []Image                       `json:"musiclabel,omitempty"`
	Albums           map[string]map[string][]Image `json:"albums,omitempty"`
}

// ID returns the id of the original query.
func (res *ImagesResult) ID() string {
	switch {
	case res.MbID != "":
		return res.MbID
	case res.TvdbID != "":
		return res.TvdbID
	case res.ImdbID != "":
		return res.ImdbID
	case res.TmdbID != "":
		return res.TmdbID
	}
	return res.FanartID
}

// Image contains information for a fanart image.
type Image struct {
	ID       string `json:"id,omitempty"`
	URL      string `json:"url,omitempty"`
	Lang     string `json:"lang,omitempty"`
	Colour   string `json:"colour,omitempty"`
	DiscType string `json:"disc_type,omitempty"`
	Likes    int    `json:"likes,omitempty"`
	Season   int    `json:"season,omitempty"`
	Disc     int    `json:"disc,omitempty"`
	Size     int    `json:"size,omitempty"`
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (img *Image) UnmarshalJSON(buf []byte) error {
	type image struct {
		ID       string `json:"id,omitempty"`
		URL      string `json:"url,omitempty"`
		Lang     string `json:"lang,omitempty"`
		Colour   string `json:"colour,omitempty"`
		DiscType string `json:"disc_type,omitempty"`
		Likes    string `json:"likes,omitempty"`
		Season   string `json:"season,omitempty"`
		Disc     string `json:"disc,omitempty"`
		Size     string `json:"size,omitempty"`
	}
	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.DisallowUnknownFields()
	v := image{}
	if err := dec.Decode(&v); err != nil {
		return err
	}
	img.ID = v.ID
	img.URL = v.URL
	img.Lang = v.Lang
	img.Colour = v.Colour
	img.DiscType = v.DiscType
	// img.Albums = v.Albums
	var err error
	if v.Likes != "" {
		if img.Likes, err = strconv.Atoi(v.Likes); err != nil {
			return fmt.Errorf("invalid likes: %w", err)
		}
	}
	if v.Season != "" && v.Season != "all" {
		if img.Season, err = strconv.Atoi(v.Season); err != nil {
			return fmt.Errorf("invalid season: %w", err)
		}
	}
	if v.Disc != "" {
		if img.Disc, err = strconv.Atoi(v.Disc); err != nil {
			return fmt.Errorf("invalid disc: %w", err)
		}
	}
	if v.Size != "" {
		if img.Size, err = strconv.Atoi(v.Size); err != nil {
			return fmt.Errorf("invalid size: %w", err)
		}
	}
	// change all urls to https
	if strings.HasPrefix(img.URL, "http://") {
		img.URL = "https://" + strings.TrimPrefix(img.URL, "http://")
	}
	return nil
}

// LatestRequest is the fanart latest request.
type LatestRequest struct {
	Type Type
}

// Latest creates a new fanart latest request.
func Latest(typ Type) *LatestRequest {
	return &LatestRequest{
		Type: typ,
	}
}

// Do executes the latest request against the client.
func (req *LatestRequest) Do(ctx context.Context, cl *Client) ([]LatestResult, error) {
	var res []LatestResult
	if err := cl.Do(ctx, req.Type.ApiType()+"/latest", &res); err != nil {
		return nil, err
	}
	return res, nil
}

// LatestResult holds latest information.
type LatestResult struct {
	FanartID    string `json:"id,omitempty"`
	TmdbID      string `json:"tmdb_id,omitempty"`
	ImdbID      string `json:"imdb_id,omitempty"`
	Name        string `json:"name,omitempty"`
	NewImages   int    `json:"new_images,omitempty"`
	TotalImages int    `json:"total_images,omitempty"`
}

// ID returns the id of the latest item.
func (res *LatestResult) ID() string {
	switch {
	case res.ImdbID != "":
		return res.ImdbID
	case res.TmdbID != "":
		return res.TmdbID
	}
	return res.FanartID
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (res *LatestResult) UnmarshalJSON(buf []byte) error {
	type latest struct {
		ID          string `json:"id,omitempty"`
		TmdbID      string `json:"tmdb_id,omitempty"`
		ImdbID      string `json:"imdb_id,omitempty"`
		Name        string `json:"name,omitempty"`
		NewImages   string `json:"new_images,omitempty"`
		TotalImages string `json:"total_images,omitempty"`
	}
	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.DisallowUnknownFields()
	v := latest{}
	if err := dec.Decode(&v); err != nil {
		return err
	}
	res.FanartID = v.ID
	res.TmdbID = v.TmdbID
	res.ImdbID = v.ImdbID
	res.Name = v.Name
	var err error
	if v.NewImages != "" {
		if res.NewImages, err = strconv.Atoi(v.NewImages); err != nil {
			return fmt.Errorf("invalid new_images: %w", err)
		}
	}
	if v.TotalImages != "" {
		if res.TotalImages, err = strconv.Atoi(v.TotalImages); err != nil {
			return fmt.Errorf("invalid total_images: %w", err)
		}
	}
	return nil
}

package netease

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/navidrome/navidrome/log"
)

const (
	searchURL      = "https://music.163.com/api/search/get/web"
	artistURL      = "http://music.163.com/api/v1/artist/%d"
	artistAlbumsURL = "http://music.163.com/api/artist/albums/%d"
	albumURL       = "http://music.163.com/api/album/%d"
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func newClient(hc httpDoer) *client {
	return &client{hc: hc}
}

type client struct {
	hc httpDoer
}

// SearchArtist searches for an artist by name and returns the artist ID
func (c *client) SearchArtist(ctx context.Context, name string) (int, error) {
	params := url.Values{}
	params.Add("csrf_token", "hlpretag=")
	params.Add("hlposttag", "")
	params.Add("s", name)
	params.Add("type", "100")
	params.Add("offset", "0")
	params.Add("total", "true")
	params.Add("limit", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return 0, fmt.Errorf("create search request: %w", err)
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	log.Trace(ctx, "Searching for artist on NetEase", "name", name)

	resp, err := c.hc.Do(req)
	if err != nil {
		return 0, fmt.Errorf("search artist: %w", err)
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return 0, fmt.Errorf("decode search response: %w", err)
	}

	if searchResp.Code != 200 || searchResp.Result.ArtistCount == 0 {
		return 0, fmt.Errorf("artist not found: %s", name)
	}

	return searchResp.Result.Artists[0].ID, nil
}

// GetArtist retrieves artist information by ID
func (c *client) GetArtist(ctx context.Context, artistID int) (*ArtistResponse, error) {
	url := fmt.Sprintf(artistURL, artistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create artist request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	log.Trace(ctx, "Getting artist info from NetEase", "id", artistID)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get artist: %w", err)
	}
	defer resp.Body.Close()

	var artistResp ArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&artistResp); err != nil {
		return nil, fmt.Errorf("decode artist response: %w", err)
	}

	if artistResp.Code != 200 {
		return nil, fmt.Errorf("artist API returned error: %d", artistResp.Code)
	}

	return &artistResp, nil
}

// GetArtistAlbums retrieves all albums for an artist
func (c *client) GetArtistAlbums(ctx context.Context, artistID int) (*ArtistAlbumsResponse, error) {
	url := fmt.Sprintf(artistAlbumsURL, artistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create artist albums request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	log.Trace(ctx, "Getting artist albums from NetEase", "id", artistID)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get artist albums: %w", err)
	}
	defer resp.Body.Close()

	var albumsResp ArtistAlbumsResponse
	if err := json.NewDecoder(resp.Body).Decode(&albumsResp); err != nil {
		return nil, fmt.Errorf("decode albums response: %w", err)
	}

	if albumsResp.Code != 200 {
		return nil, fmt.Errorf("artist albums API returned error: %d", albumsResp.Code)
	}

	return &albumsResp, nil
}

// GetAlbum retrieves album information by ID
func (c *client) GetAlbum(ctx context.Context, albumID int) (*AlbumResponse, error) {
	url := fmt.Sprintf(albumURL, albumID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create album request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	log.Trace(ctx, "Getting album info from NetEase", "id", albumID)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}
	defer resp.Body.Close()

	var albumResp AlbumResponse
	if err := json.NewDecoder(resp.Body).Decode(&albumResp); err != nil {
		return nil, fmt.Errorf("decode album response: %w", err)
	}

	if albumResp.Code != 200 {
		return nil, fmt.Errorf("album API returned error: %d", albumResp.Code)
	}

	return &albumResp, nil
}

// FindAlbumByName finds an album ID by name from a list of albums
func (c *client) FindAlbumByName(albums *ArtistAlbumsResponse, albumName string) (int, error) {
	for _, album := range albums.HotAlbums {
		if strings.EqualFold(album.Name, albumName) {
			return album.ID, nil
		}
	}
	return 0, fmt.Errorf("album not found: %s", albumName)
}

// GetAlbumByArtistAndName searches for an album by artist and album name
func (c *client) GetAlbumByArtistAndName(ctx context.Context, artistName, albumName string) (*AlbumResponse, error) {
	// First, search for the artist
	artistID, err := c.SearchArtist(ctx, artistName)
	if err != nil {
		return nil, fmt.Errorf("search artist: %w", err)
	}

	// Get all albums for the artist
	albums, err := c.GetArtistAlbums(ctx, artistID)
	if err != nil {
		return nil, fmt.Errorf("get artist albums: %w", err)
	}

	// Find the specific album
	albumID, err := c.FindAlbumByName(albums, albumName)
	if err != nil {
		return nil, fmt.Errorf("find album: %w", err)
	}

	// Get album details
	return c.GetAlbum(ctx, albumID)
}

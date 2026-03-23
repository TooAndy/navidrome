package netease

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/cache"
)

const (
	neteaseAgentName = "netease"
)

// neteaseClient is the interface for NetEase API client
type neteaseClient interface {
	SearchArtist(ctx context.Context, name string) (int, error)
	GetArtist(ctx context.Context, artistID int) (*ArtistResponse, error)
	GetArtistAlbums(ctx context.Context, artistID int) (*ArtistAlbumsResponse, error)
	GetAlbum(ctx context.Context, albumID int) (*AlbumResponse, error)
	FindAlbumByName(albums *ArtistAlbumsResponse, albumName string) (int, error)
	GetAlbumByArtistAndName(ctx context.Context, artistName, albumName string) (*AlbumResponse, error)
}

type neteaseAgent struct {
	ds          model.DataStore
	client      neteaseClient
	httpClient  httpDoer
	cache       sync.Map
}

func neteaseConstructor(ds model.DataStore) *neteaseAgent {
	if !conf.Server.NetEase.Enabled {
		return nil
	}

	n := &neteaseAgent{
		ds: ds,
	}

	hc := &http.Client{
		Timeout: consts.DefaultHttpClientTimeOut,
	}
	chc := cache.NewHTTPClient(hc, consts.DefaultHttpClientTimeOut)
	n.httpClient = chc
	n.client = newClient(chc)

	log.Info("NetEase agent enabled")

	return n
}

func (n *neteaseAgent) AgentName() string {
	return neteaseAgentName
}

// GetArtistBiography retrieves artist biography from NetEase
func (n *neteaseAgent) GetArtistBiography(ctx context.Context, id, name, mbid string) (string, error) {
	log.Debug(ctx, "Getting artist biography from NetEase", "name", name)

	artistResp, err := n.client.SearchArtist(ctx, name)
	if err != nil {
		log.Debug(ctx, "Artist not found on NetEase", "name", name, err)
		return "", agents.ErrNotFound
	}

	profile, err := n.client.GetArtist(ctx, artistResp)
	if err != nil {
		log.Debug(ctx, "Failed to get artist profile from NetEase", "name", name, err)
		return "", agents.ErrNotFound
	}

	bio := strings.TrimSpace(profile.Artist.BriefDesc)
	if bio == "" {
		bio = strings.TrimSpace(profile.Artist.Desc)
	}

	if bio == "" {
		return "", agents.ErrNotFound
	}

	return bio, nil
}

// GetArtistImages retrieves artist images from NetEase
func (n *neteaseAgent) GetArtistImages(ctx context.Context, id, name, mbid string) ([]agents.ExternalImage, error) {
	log.Debug(ctx, "Getting artist images from NetEase", "name", name)

	artistResp, err := n.client.SearchArtist(ctx, name)
	if err != nil {
		log.Debug(ctx, "Artist not found on NetEase", "name", name, err)
		return nil, agents.ErrNotFound
	}

	profile, err := n.client.GetArtist(ctx, artistResp)
	if err != nil {
		log.Debug(ctx, "Failed to get artist profile from NetEase", "name", name, err)
		return nil, agents.ErrNotFound
	}

	// NetEase returns an image URL, create multiple sizes for compatibility
	imageURL := profile.Artist.Img1V1URL
	if imageURL == "" {
		imageURL = profile.Artist.PicURL
	}

	if imageURL == "" {
		return nil, agents.ErrNotFound
	}

	images := []agents.ExternalImage{
		{URL: imageURL, Size: 160},
		{URL: imageURL, Size: 320},
		{URL: imageURL, Size: 640},
	}

	return images, nil
}

// GetAlbumInfo retrieves album information from NetEase
func (n *neteaseAgent) GetAlbumInfo(ctx context.Context, name, artist, mbid string) (*agents.AlbumInfo, error) {
	log.Debug(ctx, "Getting album info from NetEase", "album", name, "artist", artist)

	album, err := n.client.GetAlbumByArtistAndName(ctx, artist, name)
	if err != nil {
		log.Debug(ctx, "Album not found on NetEase", "album", name, "artist", artist, err)
		return nil, agents.ErrNotFound
	}

	info := &agents.AlbumInfo{
		Name:        name,
		Description: strings.TrimSpace(album.Album.Description),
	}

	if info.Description == "" {
		return info, nil
	}

	return info, nil
}

// GetAlbumImages retrieves album images from NetEase
func (n *neteaseAgent) GetAlbumImages(ctx context.Context, name, artist, mbid string) ([]agents.ExternalImage, error) {
	log.Debug(ctx, "Getting album images from NetEase", "album", name, "artist", artist)

	album, err := n.client.GetAlbumByArtistAndName(ctx, artist, name)
	if err != nil {
		log.Debug(ctx, "Album not found on NetEase", "album", name, "artist", artist, err)
		return nil, agents.ErrNotFound
	}

	imageURL := album.Album.PicURL
	if imageURL == "" {
		imageURL = album.Album.BlurPicURL
	}

	if imageURL == "" {
		return nil, agents.ErrNotFound
	}

	images := []agents.ExternalImage{
		{URL: imageURL, Size: 300},
		{URL: imageURL, Size: 600},
		{URL: imageURL, Size: 1200},
	}

	return images, nil
}

// GetArtistMBID is not supported by NetEase
func (n *neteaseAgent) GetArtistMBID(ctx context.Context, id string, name string) (string, error) {
	return "", agents.ErrNotFound
}

// GetArtistURL is not supported by NetEase
func (n *neteaseAgent) GetArtistURL(ctx context.Context, id, name, mbid string) (string, error) {
	return "", agents.ErrNotFound
}

// GetSimilarArtists is not supported by NetEase
func (n *neteaseAgent) GetSimilarArtists(ctx context.Context, id, name, mbid string, limit int) ([]agents.Artist, error) {
	return nil, agents.ErrNotFound
}

// GetArtistTopSongs is not supported by NetEase
func (n *neteaseAgent) GetArtistTopSongs(ctx context.Context, id, artistName, mbid string, count int) ([]agents.Song, error) {
	return nil, agents.ErrNotFound
}

// GetSimilarSongsByTrack is not supported by NetEase
func (n *neteaseAgent) GetSimilarSongsByTrack(ctx context.Context, id, name, artist, mbid string, count int) ([]agents.Song, error) {
	return nil, agents.ErrNotFound
}

// GetSimilarSongsByAlbum is not supported by NetEase
func (n *neteaseAgent) GetSimilarSongsByAlbum(ctx context.Context, id, name, artist, mbid string, count int) ([]agents.Song, error) {
	return nil, agents.ErrNotFound
}

// GetSimilarSongsByArtist is not supported by NetEase
func (n *neteaseAgent) GetSimilarSongsByArtist(ctx context.Context, id, name, mbid string, count int) ([]agents.Song, error) {
	return nil, agents.ErrNotFound
}

func init() {
	conf.AddHook(func() {
		agents.Register(neteaseAgentName, func(ds model.DataStore) agents.Interface {
			a := neteaseConstructor(ds)
			if a != nil {
				return a
			}
			return nil
		})
	})
}

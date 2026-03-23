package netease

import (
	"context"
	"errors"
	"strings"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/conf/configtest"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("neteaseAgent", func() {
	var ds model.DataStore
	var ctx context.Context

	BeforeEach(func() {
		ds = &tests.MockDataStore{}
		ctx = context.Background()
		DeferCleanup(configtest.SetupConfig())
		conf.Server.NetEase.Enabled = true
	})

	Describe("neteaseConstructor", func() {
		When("Agent is properly configured", func() {
			It("creates an agent", func() {
				agent := neteaseConstructor(ds)
				Expect(agent).ToNot(BeNil())
				Expect(agent.AgentName()).To(Equal("netease"))
			})
		})

		When("Agent is disabled", func() {
			It("returns nil", func() {
				conf.Server.NetEase.Enabled = false
				Expect(neteaseConstructor(ds)).To(BeNil())
			})
		})
	})

	Describe("GetArtistBiography", func() {
		var agent *neteaseAgent
		var mockClient *mockNeteaseClient

		BeforeEach(func() {
			mockClient = &mockNeteaseClient{}
			agent = neteaseConstructor(ds)
			agent.client = mockClient
		})

		It("returns the artist biography", func() {
			mockClient.artistResp = &ArtistResponse{
				Code: 200,
			}
			mockClient.artistResp.Artist.BriefDesc = "周杰伦是华语流行音乐男歌手"

			bio, err := agent.GetArtistBiography(ctx, "123", "周杰伦", "")
			Expect(err).ToNot(HaveOccurred())
			Expect(bio).To(Equal("周杰伦是华语流行音乐男歌手"))
			Expect(mockClient.searchCalled).To(BeTrue())
			Expect(mockClient.artistCalled).To(BeTrue())
		})

		It("returns ErrNotFound when artist is not found", func() {
			mockClient.searchErr = errors.New("not found")

			_, err := agent.GetArtistBiography(ctx, "123", "UnknownArtist", "")
			Expect(err).To(HaveOccurred())
			Expect(mockClient.searchCalled).To(BeTrue())
		})

		It("returns empty biography when API returns empty", func() {
			mockClient.artistResp = &ArtistResponse{
				Code: 200,
			}

			_, err := agent.GetArtistBiography(ctx, "123", "周杰伦", "")
			Expect(err).To(MatchError(agents.ErrNotFound))
		})
	})

	Describe("GetArtistImages", func() {
		var agent *neteaseAgent
		var mockClient *mockNeteaseClient

		BeforeEach(func() {
			mockClient = &mockNeteaseClient{}
			agent = neteaseConstructor(ds)
			agent.client = mockClient
		})

		It("returns artist images", func() {
			mockClient.artistResp = &ArtistResponse{
				Code: 200,
			}
			mockClient.artistResp.Artist.Img1V1URL = "http://example.com/image.jpg"

			images, err := agent.GetArtistImages(ctx, "123", "周杰伦", "")
			Expect(err).ToNot(HaveOccurred())
			Expect(images).To(HaveLen(3))
			Expect(images[0].Size).To(Equal(160))
			Expect(images[0].URL).To(Equal("http://example.com/image.jpg"))
		})

		It("returns ErrNotFound when artist is not found", func() {
			mockClient.searchErr = errors.New("not found")

			_, err := agent.GetArtistImages(ctx, "123", "UnknownArtist", "")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetAlbumInfo", func() {
		var agent *neteaseAgent
		var mockClient *mockNeteaseClient

		BeforeEach(func() {
			mockClient = &mockNeteaseClient{}
			agent = neteaseConstructor(ds)
			agent.client = mockClient
		})

		It("returns album information", func() {
			mockClient.albumResp = &AlbumResponse{
				Code: 200,
			}
			mockClient.albumResp.Album.Description = "叶惠美是周杰伦的第四张专辑"

			info, err := agent.GetAlbumInfo(ctx, "叶惠美", "周杰伦", "")
			Expect(err).ToNot(HaveOccurred())
			Expect(info.Name).To(Equal("叶惠美"))
			Expect(info.Description).To(Equal("叶惠美是周杰伦的第四张专辑"))
		})

		It("returns ErrNotFound when album is not found", func() {
			mockClient.searchErr = errors.New("not found")

			_, err := agent.GetAlbumInfo(ctx, "UnknownAlbum", "UnknownArtist", "")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetAlbumImages", func() {
		var agent *neteaseAgent
		var mockClient *mockNeteaseClient

		BeforeEach(func() {
			mockClient = &mockNeteaseClient{}
			agent = neteaseConstructor(ds)
			agent.client = mockClient
		})

		It("returns album images", func() {
			mockClient.albumResp = &AlbumResponse{
				Code: 200,
			}
			mockClient.albumResp.Album.PicURL = "http://example.com/album.jpg"

			images, err := agent.GetAlbumImages(ctx, "叶惠美", "周杰伦", "")
			Expect(err).ToNot(HaveOccurred())
			Expect(images).To(HaveLen(3))
			Expect(images[0].Size).To(Equal(300))
			Expect(images[0].URL).To(Equal("http://example.com/album.jpg"))
		})

		It("returns ErrNotFound when album is not found", func() {
			mockClient.searchErr = errors.New("not found")

			_, err := agent.GetAlbumImages(ctx, "UnknownAlbum", "UnknownArtist", "")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Unsupported operations", func() {
		var agent *neteaseAgent

		BeforeEach(func() {
			agent = neteaseConstructor(ds)
		})

		It("GetArtistMBID returns ErrNotFound", func() {
			_, err := agent.GetArtistMBID(ctx, "123", "周杰伦")
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetArtistURL returns ErrNotFound", func() {
			_, err := agent.GetArtistURL(ctx, "123", "周杰伦", "")
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetSimilarArtists returns ErrNotFound", func() {
			_, err := agent.GetSimilarArtists(ctx, "123", "周杰伦", "", 10)
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetArtistTopSongs returns ErrNotFound", func() {
			_, err := agent.GetArtistTopSongs(ctx, "123", "周杰伦", "", 10)
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetSimilarSongsByTrack returns ErrNotFound", func() {
			_, err := agent.GetSimilarSongsByTrack(ctx, "123", "晴天", "周杰伦", "", 10)
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetSimilarSongsByAlbum returns ErrNotFound", func() {
			_, err := agent.GetSimilarSongsByAlbum(ctx, "123", "叶惠美", "周杰伦", "", 10)
			Expect(err).To(MatchError(agents.ErrNotFound))
		})

		It("GetSimilarSongsByArtist returns ErrNotFound", func() {
			_, err := agent.GetSimilarSongsByArtist(ctx, "123", "周杰伦", "", 10)
			Expect(err).To(MatchError(agents.ErrNotFound))
		})
	})
})

// mockNeteaseClient is a mock implementation of the netease client for testing
type mockNeteaseClient struct {
	searchCalled bool
	artistCalled bool
	albumsCalled bool
	albumCalled  bool

	searchErr     error
	artistResp    *ArtistResponse
	artistErr     error
	albumsResp    *ArtistAlbumsResponse
	albumsErr     error
	albumResp     *AlbumResponse
	albumErr      error
	findAlbumID   int
	findAlbumErr  error
}

func (m *mockNeteaseClient) SearchArtist(ctx context.Context, name string) (int, error) {
	m.searchCalled = true
	if m.searchErr != nil {
		return 0, m.searchErr
	}
	return 6452, nil
}

func (m *mockNeteaseClient) GetArtist(ctx context.Context, artistID int) (*ArtistResponse, error) {
	m.artistCalled = true
	if m.artistErr != nil {
		return nil, m.artistErr
	}
	return m.artistResp, nil
}

func (m *mockNeteaseClient) GetArtistAlbums(ctx context.Context, artistID int) (*ArtistAlbumsResponse, error) {
	m.albumsCalled = true
	if m.albumsErr != nil {
		return nil, m.albumsErr
	}
	if m.albumsResp == nil {
		m.albumsResp = &ArtistAlbumsResponse{
			Code: 200,
			HotAlbums: []struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				PicURL string `json:"picUrl"`
			}{{ID: 12345, Name: "叶惠美"}},
		}
	}
	return m.albumsResp, nil
}

func (m *mockNeteaseClient) GetAlbum(ctx context.Context, albumID int) (*AlbumResponse, error) {
	m.albumCalled = true
	if m.albumErr != nil {
		return nil, m.albumErr
	}
	return m.albumResp, nil
}

func (m *mockNeteaseClient) FindAlbumByName(albums *ArtistAlbumsResponse, albumName string) (int, error) {
	if m.findAlbumErr != nil {
		return 0, m.findAlbumErr
	}
	if m.findAlbumID != 0 {
		return m.findAlbumID, nil
	}
	for _, album := range albums.HotAlbums {
		if strings.EqualFold(album.Name, albumName) {
			return album.ID, nil
		}
	}
	return 0, errors.New("album not found")
}

func (m *mockNeteaseClient) GetAlbumByArtistAndName(ctx context.Context, artistName, albumName string) (*AlbumResponse, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	if m.albumResp == nil {
		return nil, errors.New("album not found")
	}
	return m.albumResp, nil
}

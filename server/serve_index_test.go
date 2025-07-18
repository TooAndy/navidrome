package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/conf/configtest"
	"github.com/navidrome/navidrome/conf/mime"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("serveIndex", func() {
	var ds model.DataStore
	mockUser := &mockedUserRepo{}
	fs := os.DirFS("tests/fixtures")

	BeforeEach(func() {
		ds = &tests.MockDataStore{MockedUser: mockUser}
		DeferCleanup(configtest.SetupConfig())
	})

	It("adds app_config to index.html", func() {
		r := httptest.NewRequest("GET", "/index.html", nil)
		w := httptest.NewRecorder()

		serveIndex(ds, fs, nil)(w, r)

		Expect(w.Code).To(Equal(200))
		config := extractAppConfig(w.Body.String())
		Expect(config).To(BeAssignableToTypeOf(map[string]any{}))
	})

	It("sets firstTime = true when User table is empty", func() {
		mockUser.empty = true
		r := httptest.NewRequest("GET", "/index.html", nil)
		w := httptest.NewRecorder()

		serveIndex(ds, fs, nil)(w, r)

		config := extractAppConfig(w.Body.String())
		Expect(config).To(HaveKeyWithValue("firstTime", true))
	})

	It("sets firstTime = false when User table is not empty", func() {
		mockUser.empty = false
		r := httptest.NewRequest("GET", "/index.html", nil)
		w := httptest.NewRecorder()

		serveIndex(ds, fs, nil)(w, r)

		config := extractAppConfig(w.Body.String())
		Expect(config).To(HaveKeyWithValue("firstTime", false))
	})

	DescribeTable("sets configuration values",
		func(configSetter func(), configKey string, expectedValue any) {
			configSetter()
			r := httptest.NewRequest("GET", "/index.html", nil)
			w := httptest.NewRecorder()

			serveIndex(ds, fs, nil)(w, r)

			config := extractAppConfig(w.Body.String())
			Expect(config).To(HaveKeyWithValue(configKey, expectedValue))
		},
		Entry("baseURL", func() { conf.Server.BasePath = "base_url_test" }, "baseURL", "base_url_test"),
		Entry("welcomeMessage", func() { conf.Server.UIWelcomeMessage = "Hello" }, "welcomeMessage", "Hello"),
		Entry("maxSidebarPlaylists", func() { conf.Server.MaxSidebarPlaylists = 42 }, "maxSidebarPlaylists", float64(42)),
		Entry("enableTranscodingConfig", func() { conf.Server.EnableTranscodingConfig = true }, "enableTranscodingConfig", true),
		Entry("enableDownloads", func() { conf.Server.EnableDownloads = true }, "enableDownloads", true),
		Entry("enableFavourites", func() { conf.Server.EnableFavourites = true }, "enableFavourites", true),
		Entry("enableStarRating", func() { conf.Server.EnableStarRating = true }, "enableStarRating", true),
		Entry("defaultTheme", func() { conf.Server.DefaultTheme = "Light" }, "defaultTheme", "Light"),
		Entry("defaultLanguage", func() { conf.Server.DefaultLanguage = "pt" }, "defaultLanguage", "pt"),
		Entry("defaultUIVolume", func() { conf.Server.DefaultUIVolume = 45 }, "defaultUIVolume", float64(45)),
		Entry("enableCoverAnimation", func() { conf.Server.EnableCoverAnimation = true }, "enableCoverAnimation", true),
		Entry("enableNowPlaying", func() { conf.Server.EnableNowPlaying = true }, "enableNowPlaying", true),
		Entry("gaTrackingId", func() { conf.Server.GATrackingID = "UA-12345" }, "gaTrackingId", "UA-12345"),
		Entry("defaultDownloadableShare", func() { conf.Server.DefaultDownloadableShare = true }, "defaultDownloadableShare", true),
		Entry("devSidebarPlaylists", func() { conf.Server.DevSidebarPlaylists = true }, "devSidebarPlaylists", true),
		Entry("lastFMEnabled", func() { conf.Server.LastFM.Enabled = true }, "lastFMEnabled", true),
		Entry("devShowArtistPage", func() { conf.Server.DevShowArtistPage = true }, "devShowArtistPage", true),
		Entry("devUIShowConfig", func() { conf.Server.DevUIShowConfig = true }, "devUIShowConfig", true),
		Entry("listenBrainzEnabled", func() { conf.Server.ListenBrainz.Enabled = true }, "listenBrainzEnabled", true),
		Entry("enableReplayGain", func() { conf.Server.EnableReplayGain = true }, "enableReplayGain", true),
		Entry("enableExternalServices", func() { conf.Server.EnableExternalServices = true }, "enableExternalServices", true),
		Entry("devActivityPanel", func() { conf.Server.DevActivityPanel = true }, "devActivityPanel", true),
		Entry("shareURL", func() { conf.Server.ShareURL = "https://share.example.com" }, "shareURL", "https://share.example.com"),
		Entry("enableInspect", func() { conf.Server.Inspect.Enabled = true }, "enableInspect", true),
		Entry("defaultDownsamplingFormat", func() { conf.Server.DefaultDownsamplingFormat = "mp3" }, "defaultDownsamplingFormat", "mp3"),
		Entry("enableUserEditing", func() { conf.Server.EnableUserEditing = false }, "enableUserEditing", false),
		Entry("enableSharing", func() { conf.Server.EnableSharing = true }, "enableSharing", true),
		Entry("devNewEventStream", func() { conf.Server.DevNewEventStream = true }, "devNewEventStream", true),
	)

	DescribeTable("sets other UI configuration values",
		func(configKey string, expectedValueFunc func() any) {
			r := httptest.NewRequest("GET", "/index.html", nil)
			w := httptest.NewRecorder()

			serveIndex(ds, fs, nil)(w, r)

			config := extractAppConfig(w.Body.String())
			Expect(config).To(HaveKeyWithValue(configKey, expectedValueFunc()))
		},
		Entry("version", "version", func() any { return consts.Version }),
		Entry("variousArtistsId", "variousArtistsId", func() any { return consts.VariousArtistsID }),
		Entry("losslessFormats", "losslessFormats", func() any {
			return strings.ToUpper(strings.Join(mime.LosslessFormats, ","))
		}),
		Entry("separator", "separator", func() any { return string(os.PathSeparator) }),
	)

	Describe("loginBackgroundURL", func() {
		Context("empty BaseURL", func() {
			BeforeEach(func() {
				conf.Server.BasePath = "/"
			})
			When("it is the default URL", func() {
				It("points to the default URL", func() {
					conf.Server.UILoginBackgroundURL = consts.DefaultUILoginBackgroundURL
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", consts.DefaultUILoginBackgroundURL))
				})
			})
			When("it is the default offline URL", func() {
				It("points to the offline URL", func() {
					conf.Server.UILoginBackgroundURL = consts.DefaultUILoginBackgroundURLOffline
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", consts.DefaultUILoginBackgroundURLOffline))
				})
			})
			When("it is a custom URL", func() {
				It("points to the offline URL", func() {
					conf.Server.UILoginBackgroundURL = "https://example.com/images/1.jpg"
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", "https://example.com/images/1.jpg"))
				})
			})
		})
		Context("with a BaseURL", func() {
			BeforeEach(func() {
				conf.Server.BasePath = "/music"
			})
			When("it is the default URL", func() {
				It("points to the default URL with BaseURL prefix", func() {
					conf.Server.UILoginBackgroundURL = consts.DefaultUILoginBackgroundURL
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", "/music"+consts.DefaultUILoginBackgroundURL))
				})
			})
			When("it is the default offline URL", func() {
				It("points to the offline URL", func() {
					conf.Server.UILoginBackgroundURL = consts.DefaultUILoginBackgroundURLOffline
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", consts.DefaultUILoginBackgroundURLOffline))
				})
			})
			When("it is a custom URL", func() {
				It("points to the offline URL", func() {
					conf.Server.UILoginBackgroundURL = "https://example.com/images/1.jpg"
					r := httptest.NewRequest("GET", "/index.html", nil)
					w := httptest.NewRecorder()

					serveIndex(ds, fs, nil)(w, r)

					config := extractAppConfig(w.Body.String())
					Expect(config).To(HaveKeyWithValue("loginBackgroundURL", "https://example.com/images/1.jpg"))
				})
			})
		})
	})
})

var _ = Describe("addShareData", func() {
	var (
		r         *http.Request
		data      map[string]any
		shareInfo *model.Share
	)

	BeforeEach(func() {
		data = make(map[string]any)
		r = httptest.NewRequest("GET", "/", nil)
	})

	Context("when shareInfo is nil or has an empty ID", func() {
		It("should not modify data", func() {
			addShareData(r, data, nil)
			Expect(data).To(BeEmpty())

			shareInfo = &model.Share{}
			addShareData(r, data, shareInfo)
			Expect(data).To(BeEmpty())
		})
	})

	Context("when shareInfo is not nil and has a non-empty ID", func() {
		BeforeEach(func() {
			shareInfo = &model.Share{
				ID:           "testID",
				Description:  "Test description",
				Downloadable: true,
				Tracks: []model.MediaFile{
					{
						ID:        "track1",
						Title:     "Track 1",
						Artist:    "Artist 1",
						Album:     "Album 1",
						Duration:  100,
						UpdatedAt: time.Date(2023, time.Month(3), 27, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:        "track2",
						Title:     "Track 2",
						Artist:    "Artist 2",
						Album:     "Album 2",
						Duration:  200,
						UpdatedAt: time.Date(2023, time.Month(3), 26, 0, 0, 0, 0, time.UTC),
					},
				},
				Contents: "Test contents",
				URL:      "https://example.com/share/testID",
				ImageURL: "https://example.com/share/testID/image",
			}
		})

		It("should populate data with shareInfo data", func() {
			addShareData(r, data, shareInfo)

			Expect(data["ShareDescription"]).To(Equal(shareInfo.Description))
			Expect(data["ShareURL"]).To(Equal(shareInfo.URL))
			Expect(data["ShareImageURL"]).To(Equal(shareInfo.ImageURL))

			var shareData shareData
			err := json.Unmarshal([]byte(data["ShareInfo"].(string)), &shareData)
			Expect(err).NotTo(HaveOccurred())
			Expect(shareData.ID).To(Equal(shareInfo.ID))
			Expect(shareData.Description).To(Equal(shareInfo.Description))
			Expect(shareData.Downloadable).To(Equal(shareInfo.Downloadable))

			Expect(shareData.Tracks).To(HaveLen(len(shareInfo.Tracks)))
			for i, track := range shareData.Tracks {
				Expect(track.ID).To(Equal(shareInfo.Tracks[i].ID))
				Expect(track.Title).To(Equal(shareInfo.Tracks[i].Title))
				Expect(track.Artist).To(Equal(shareInfo.Tracks[i].Artist))
				Expect(track.Album).To(Equal(shareInfo.Tracks[i].Album))
				Expect(track.Duration).To(Equal(shareInfo.Tracks[i].Duration))
				Expect(track.UpdatedAt).To(Equal(shareInfo.Tracks[i].UpdatedAt))
			}
		})

		Context("when shareInfo has an empty description", func() {
			BeforeEach(func() {
				shareInfo.Description = ""
			})

			It("should use shareInfo.Contents as ShareDescription", func() {
				addShareData(r, data, shareInfo)
				Expect(data["ShareDescription"]).To(Equal(shareInfo.Contents))
			})
		})
	})
})

var appConfigRegex = regexp.MustCompile(`(?m)window.__APP_CONFIG__=(.*);</script>`)

func extractAppConfig(body string) map[string]any {
	config := make(map[string]any)
	match := appConfigRegex.FindStringSubmatch(body)
	if match == nil {
		return config
	}
	str, err := strconv.Unquote(match[1])
	if err != nil {
		panic(fmt.Sprintf("%s: %s", match[1], err))
	}
	if err := json.Unmarshal([]byte(str), &config); err != nil {
		panic(err)
	}
	return config
}

type mockedUserRepo struct {
	model.UserRepository
	empty bool
}

func (u *mockedUserRepo) CountAll(...model.QueryOptions) (int64, error) {
	if u.empty {
		return 0, nil
	}
	return 1, nil
}

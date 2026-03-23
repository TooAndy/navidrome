package netease

// SearchResponse is the response from the search API
type SearchResponse struct {
	Result struct {
		ArtistCount int `json:"artistCount"`
		Artists     []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"result"`
	Code int `json:"code"`
}

// ArtistResponse is the response from the artist API
type ArtistResponse struct {
	Code int `json:"code"`
	Artist struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		PicURL     string `json:"picUrl"`
		Img1V1URL  string `json:"img1v1Url"`
		BriefDesc  string `json:"briefDesc"`
		Desc       string `json:"desc"`
	} `json:"artist"`
}

// ArtistAlbumsResponse is the response from the artist albums API
type ArtistAlbumsResponse struct {
	Code int `json:"code"`
	HotAlbums []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	} `json:"hotAlbums"`
}

// AlbumResponse is the response from the album API
type AlbumResponse struct {
	Code   int `json:"code"`
	Album struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		PicURL      string `json:"picUrl"`
		BlurPicURL  string `json:"blurPicUrl"`
		Description string `json:"description"`
	} `json:"album"`
}

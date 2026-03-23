//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

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

type ArtistResponse struct {
	Code   int `json:"code"`
	Artist struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		PicURL    string `json:"picUrl"`
		BriefDesc string `json:"briefDesc"`
	} `json:"artist"`
}

type ArtistAlbumsResponse struct {
	Code      int `json:"code"`
	HotAlbums []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	} `json:"hotAlbums"`
}

type AlbumResponse struct {
	Code  int `json:"code"`
	Album struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		PicURL      string `json:"picUrl"`
		Description string `json:"description"`
	} `json:"album"`
}

func main() {
	ctx := context.Background()
	hc := &http.Client{Timeout: 10 * time.Second}

	// 测试艺术家：许嵩
	artistName := "许嵩"

	// 1. 搜索艺术家
	fmt.Printf("=== 搜索艺术家: %s ===\n", artistName)
	artistID, err := searchArtist(ctx, hc, artistName)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}
	fmt.Printf("找到艺术家 ID: %d\n\n", artistID)

	// 2. 获取艺术家详情
	fmt.Printf("=== 获取艺术家详情 ===\n")
	artistInfo, err := getArtist(ctx, hc, artistID)
	if err != nil {
		log.Fatalf("获取艺术家信息失败: %v", err)
	}
	fmt.Printf("名称: %s\n", artistInfo.Artist.Name)
	fmt.Printf("简介: %s\n", artistInfo.Artist.BriefDesc)
	fmt.Printf("图片: %s\n\n", artistInfo.Artist.PicURL)

	// 3. 获取艺术家专辑
	fmt.Printf("=== 获取艺术家专辑 ===\n")
	albums, err := getArtistAlbums(ctx, hc, artistID)
	if err != nil {
		log.Fatalf("获取专辑失败: %v", err)
	}
	fmt.Printf("找到 %d 张专辑\n", len(albums.HotAlbums))
	for i, album := range albums.HotAlbums {
		if i < 5 {
			fmt.Printf("  - %s (ID: %d)\n", album.Name, album.ID)
		}
	}
	fmt.Println()

	// 4. 获取第一张专辑详情
	if len(albums.HotAlbums) > 0 {
		albumID := albums.HotAlbums[0].ID
		fmt.Printf("=== 获取专辑详情 ===\n")
		albumInfo, err := getAlbum(ctx, hc, albumID)
		if err != nil {
			log.Fatalf("获取专辑详情失败: %v", err)
		}
		fmt.Printf("专辑名称: %s\n", albumInfo.Album.Name)
		fmt.Printf("描述: %s\n", albumInfo.Album.Description)
		fmt.Printf("封面: %s\n", albumInfo.Album.PicURL)
	}

	fmt.Println("\n=== 测试完成 ===")
}

func searchArtist(ctx context.Context, hc *http.Client, name string) (int, error) {
	params := url.Values{}
	params.Add("csrf_token", "hlpretag=")
	params.Add("hlposttag", "")
	params.Add("s", name)
	params.Add("type", "100")
	params.Add("offset", "0")
	params.Add("total", "true")
	params.Add("limit", "1")

	req, err := http.NewRequestWithContext(ctx, "GET", "https://music.163.com/api/search/get/web?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	resp, err := hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if result.Code != 200 || result.Result.ArtistCount == 0 {
		return 0, fmt.Errorf("未找到艺术家: %s", name)
	}

	return result.Result.Artists[0].ID, nil
}

func getArtist(ctx context.Context, hc *http.Client, id int) (*ArtistResponse, error) {
	url := fmt.Sprintf("http://music.163.com/api/v1/artist/%d", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ArtistResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API 返回错误: %d", result.Code)
	}

	return &result, nil
}

func getArtistAlbums(ctx context.Context, hc *http.Client, id int) (*ArtistAlbumsResponse, error) {
	url := fmt.Sprintf("http://music.163.com/api/artist/albums/%d?offset=0&total=true&limit=10", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ArtistAlbumsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API 返回错误: %d", result.Code)
	}

	return &result, nil
}

func getAlbum(ctx context.Context, hc *http.Client, id int) (*AlbumResponse, error) {
	url := fmt.Sprintf("http://music.163.com/api/album/%d", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0")
	req.Header.Set("Referer", "https://music.163.com")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AlbumResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API 返回错误: %d", result.Code)
	}

	return &result, nil
}

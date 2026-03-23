//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	// 创建 HTTP client
	hc := &http.Client{Timeout: 10 * time.Second}

	// 测试 1: 搜索艺术家
	fmt.Println("=== 测试 1: 搜索艺术家 '周杰伦' ===")
	artistID, err := searchArtist(ctx, hc, "周杰伦")
	if err != nil {
		log.Printf("搜索失败: %v", err)
		return
	}
	log.Printf("找到艺术家 ID: %d", artistID)

	// 测试 2: 获取艺术家信息
	fmt.Println("\n=== 测试 2: 获取艺术家信息 ===")
	artistInfo, err := getArtist(ctx, hc, artistID)
	if err != nil {
		log.Printf("获取艺术家信息失败: %v", err)
		return
	}
	log.Printf("艺术家名称: %s", artistInfo.Artist.Name)
	log.Printf("简介: %s", artistInfo.Artist.BriefDesc)
	log.Printf("图片: %s", artistInfo.Artist.PicURL)

	// 测试 3: 获取艺术家专辑
	fmt.Println("\n=== 测试 3: 获取艺术家专辑 ===")
	albums, err := getArtistAlbums(ctx, hc, artistID)
	if err != nil {
		log.Printf("获取专辑失败: %v", err)
		return
	}
	log.Printf("找到 %d 张专辑", len(albums.HotAlbums))
	for i, album := range albums.HotAlbums {
		if i < 5 {
			log.Printf("  - %s (ID: %d)", album.Name, album.ID)
		}
	}

	// 测试 4: 获取专辑详情
	if len(albums.HotAlbums) > 0 {
		albumID := albums.HotAlbums[0].ID
		fmt.Println("\n=== 测试 4: 获取专辑详情 ===")
		albumInfo, err := getAlbum(ctx, hc, albumID)
		if err != nil {
			log.Printf("获取专辑详情失败: %v", err)
			return
		}
		log.Printf("专辑名称: %s", albumInfo.Album.Name)
		log.Printf("描述: %s", albumInfo.Album.Description)
		log.Printf("封面: %s", albumInfo.Album.PicURL)
	}

	fmt.Println("\n=== 测试完成 ===")
}

// 简化的 API 调用函数
func searchArtist(ctx context.Context, hc *http.Client, name string) (int, error) {
	// 这里需要实际的 API 调用代码
	return 6452, nil
}

func getArtist(ctx context.Context, hc *http.Client, id int) (*ArtistResponse, error) {
	return nil, nil
}

func getArtistAlbums(ctx context.Context, hc *http.Client, id int) (*ArtistAlbumsResponse, error) {
	return nil, nil
}

func getAlbum(ctx context.Context, hc *http.Client, id int) (*AlbumResponse, error) {
	return nil, nil
}

// 响应结构
type ArtistResponse struct {
	Code int `json:"code"`
	Artist struct {
		ID        int `json:"id"`
		Name      string `json:"name"`
		PicURL    string `json:"picUrl"`
		BriefDesc string `json:"briefDesc"`
	} `json:"artist"`
}

type ArtistAlbumsResponse struct {
	Code int `json:"code"`
	HotAlbums []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	} `json:"hotAlbums"`
}

type AlbumResponse struct {
	Code int `json:"code"`
	Album struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		PicURL      string `json:"picUrl"`
		Description string `json:"description"`
	} `json:"album"`
}

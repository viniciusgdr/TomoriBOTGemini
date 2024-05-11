package tiktokServices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	pocbiServices "tomoribot-geminiai-version/src/services/pocbi"
)
type TikTok struct {
	Title string `json:"title"`
	Wmplay string `json:"wmplay"`
	Play string `json:"play"`
	Hdplay string `json:"hdplay"`
	Music string `json:"music"`
}

func TikTokDownloader(url string) (TikTok, error) {
	body := fmt.Sprintf("url=%s&count=12&cursor=0&web=1&hd=1", url)
	reader := bytes.NewBufferString(body)
	createRequest, err := http.NewRequest("POST", "https://www.tikwm.com/api/", reader)
	if err != nil {
		return TikTok{}, err
	}
	createRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	createRequest.Header.Add("Referer", "https://www.tikwm.com/")
	createRequest.Header.Add("Referrer-Policy", "strict-origin-when-cross-origin")
	createRequest.Header.Add("Sec-Fetch-Dest", "empty")
	createRequest.Header.Add("Sec-Fetch-Mode", "cors")
	createRequest.Header.Add("Sec-Fetch-Site", "same-origin")
	createRequest.Header.Add("Sec-Gpc", "1")
	createRequest.Header.Add("X-Requested-With", "XMLHttpRequest")
	createRequest.Header.Add("Cookie", "current_language=en; _ga_5370HT04Z3=GS1.1.1662231100.1.0.1662231100.0.0.0; _ga=GA1.1.949326586.1662231100; _gcl_au=1.1.1040210256.1662231101; __gads=ID=9d84e2b26498743d-22c7a756247d0099:T=1662231101:RT=1662231101:S=ALNI_MacNE4oM0cusVen0dOMUp8XpOBclA; __gpi=UID=00000931e5f18265:T=1662231101:RT=1662231101:S=ALNI_MYqb6E4N70Ey_wbGUYKayD_hbdiIg")
	createRequest.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	createRequest.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	createRequest.Header.Add("Sec-Ch-Ua-Platform", "\"Windows\"")
	createRequest.Header.Add("Accept", "application/json, text/javascript, /; q=0.01")
	createRequest.Header.Add("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	createRequest.Header.Add("Cache-Control", "no-cache")
	createRequest.Header.Add("Pragma", "no-cache")
	createRequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(createRequest)
	if err != nil {
		return TikTok{}, err
	}

	defer resp.Body.Close()

	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return TikTok{}, err
	}
	if response.(map[string]interface{})["data"] == nil {
		return TikTok{}, fmt.Errorf("video not found")
	}
	data := response.(map[string]interface{})["data"].(map[string]interface{})
	if response.(map[string]interface{})["code"] == -1 || response.(map[string]interface{})["code"] == "-1" {
		pocbi, err := pocbiServices.PocbiDownloader(url)
		if err != nil {
			return TikTok{}, err
		}
		var watemark string
		var noWatemark string
		var mp3 string
		for _, media := range pocbi.Medias {
			if media.Quality == "watermark" {
				watemark = media.URL
			}
			if media.Quality == "hd" {
				noWatemark = media.URL
			}
			if media.Quality == "mp3" {
				mp3 = media.URL
			}
		}
		if watemark == "" && noWatemark == "" && mp3 == "" {
			return TikTok{}, fmt.Errorf("video not found")
		}
		return TikTok{
			Title: pocbi.Title,
			Wmplay: watemark,
			Play: noWatemark,
			Hdplay: noWatemark,
			Music: mp3,
		}, nil
	}
	if response.(map[string]interface{})["msg"] != "success" {
		return TikTok{}, fmt.Errorf("video not found")
	}

	var tiktok TikTok
	tiktok.Title = data["title"].(string)
	tiktok.Play = "https://www.tikwm.com" + data["play"].(string)
	tiktok.Wmplay = "https://www.tikwm.com" + data["wmplay"].(string)
	tiktok.Hdplay = "https://www.tikwm.com" + data["hdplay"].(string)
	tiktok.Music = "https://www.tikwm.com" + data["music"].(string)

	return tiktok, nil
}
package pocbiServices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Pocbi struct {
	Title     string       `json:"title"`
	URL       string       `json:"url"`
	Thumbnail string       `json:"thumbnail"`
	Source    string       `json:"source"`
	Medias    []MediaPocbi `json:"medias"`
}

type MediaPocbi struct {
	Extension string  `json:"extension"`
	Quality   string  `json:"quality"`
	Size      float64 `json:"size"`
	URL       string  `json:"url"`
}

func GetPocbiToken() (string, error) {
	response, err := http.Get("https://pocbi.com/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Parse the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}
	token, exists := document.Find("#token").Attr("value")
	if !exists {
		return "", fmt.Errorf("Token not found")
	}
	return token, nil
}
func PocbiDownloader(url string) (Pocbi, error) {
	token, err := GetPocbiToken()
	if err != nil {
		return Pocbi{}, err
	}
	body := fmt.Sprintf("url=%s&token=%s", url, token)
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "https://pocbi.com/wp-json/aio-dl/video-data/", reader)
	if err != nil {
		return Pocbi{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.8")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Brave\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Cookie", "pll_language=pt; __cf_bm=pVFXrP5Xts99rLTLUm4RjkE_gwQUiVxhjgxnlQ.ePTE-1688221280-0-ASeqb2sWgGBrHgxKACHQLwidCn2PdyAg8RZKWqk1vOwjIG7hcOyolPoWEQBMsdjJug==")
	req.Header.Set("Referer", "https://pocbi.com/")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Pocbi{}, err
	}

	var pocbi Pocbi
	if response.(map[string]interface{})["title"] != nil {
		pocbi.Title = response.(map[string]interface{})["title"].(string)
	}
	if response.(map[string]interface{})["url"] != nil {
		pocbi.URL = response.(map[string]interface{})["url"].(string)
	}
	if response.(map[string]interface{})["thumbnail"] != nil {
		pocbi.Thumbnail = response.(map[string]interface{})["thumbnail"].(string)
	}
	if response.(map[string]interface{})["source"] != nil {
		pocbi.Source = response.(map[string]interface{})["source"].(string)
	}
	if response.(map[string]interface{})["medias"] == nil {
		return Pocbi{}, errors.New("media not found")
	}
	var medias []MediaPocbi
	for _, media := range response.(map[string]interface{})["medias"].([]interface{}) {
		var mediaPocbi MediaPocbi
		mediaPocbi.Extension = media.(map[string]interface{})["extension"].(string)
		mediaPocbi.Quality = media.(map[string]interface{})["quality"].(string)
		mediaPocbi.Size = media.(map[string]interface{})["size"].(float64)
		mediaPocbi.URL = media.(map[string]interface{})["url"].(string)
		medias = append(medias, mediaPocbi)
	}
	pocbi.Medias = medias

	return pocbi, nil
}

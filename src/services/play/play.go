package playServices

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type IYTSearch struct {
	VideoID   string `json:"videoId"`
	Title     string `json:"title"`
	Duration  string `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Embed     string `json:"embed"`
}
type Quality struct {
	Quality  string
	Url      string
	MimeType string
}

func Search(query string) ([]IYTSearch, error) {
	urlStr := "https://music.youtube.com/youtubei/v1/search?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false"
	requestBody := []byte(fmt.Sprintf(`{context: {client: {hl: 'pt', gl: 'BR', clientName: 'WEB', clientVersion: '2.20230602.01.00' }}, query: "%s"}`, query))
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data.(map[string]interface{})["contents"] == nil {
		return nil, errors.New("not found")
	}
	if data.(map[string]interface{})["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"] == nil {
		return nil, errors.New("not found")
	}

	resultVideos := data.(map[string]interface{})["contents"].(map[string]interface{})["twoColumnSearchResultsRenderer"].(map[string]interface{})["primaryContents"].(map[string]interface{})["sectionListRenderer"].(map[string]interface{})["contents"].([]interface{})
	var goToResults []interface{}
	switch {
	case len(resultVideos) == 0:
		goToResults = resultVideos[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})
	case len(resultVideos) > 0:
		if resultVideos[1].(map[string]interface{})["itemSectionRenderer"] != nil {
			goToResults = append(goToResults, resultVideos[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{}))
			goToResults = append(goToResults, resultVideos[1].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{}))
		} else {
			goToResults = resultVideos[0].(map[string]interface{})["itemSectionRenderer"].(map[string]interface{})["contents"].([]interface{})
		}
	}
	var videos []IYTSearch
	for _, item := range goToResults {
		if item.(map[string]interface{})["videoRenderer"] != nil {
			videoRenderer := item.(map[string]interface{})["videoRenderer"]
			videoId := videoRenderer.(map[string]interface{})["videoId"].(string)
			if videoRenderer.(map[string]interface{})["title"] == nil {
				continue
			}
			if videoRenderer.(map[string]interface{})["title"].(map[string]interface{})["runs"] == nil {
				continue
			}
			if videoRenderer.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"] == nil {
				continue
			}
			title := videoRenderer.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string)
			duration := ""
			if videoRenderer.(map[string]interface{})["lengthText"] != nil {
				duration = videoRenderer.(map[string]interface{})["lengthText"].(map[string]interface{})["simpleText"].(string)
			}
			video := IYTSearch{
				VideoID:   videoId,
				Title:     title,
				Duration:  duration,
				Thumbnail: "https://i.ytimg.com/vi/" + videoId + "/hqdefault.jpg",
				URL:       "https://www.youtube.com/watch?v=" + videoId,
				Embed:     "https://www.youtube.com/embed/" + videoId,
			}
			videos = append(videos, video)
		}
	}
	return videos, nil
}

type Info struct {
	Title       string
	Description string
	Thumbnail   string
	Author      string
	ViewCount   string
	Duration    string
}
type Streamings struct {
	Video []Quality
	Audio []Quality
}

func (info Streamings) GetHighAudio() (Quality, error) {
	var audioQualityHigh *Quality
	var audioQualityMedium *Quality
	var audioQualityLow *Quality
	for _, audio := range info.Audio {
		if audio.Quality == "AUDIO_QUALITY_MEDIUM" {
			audioQualityMedium = &audio
		} else if audio.Quality == "AUDIO_QUALITY_LOW" {
			audioQualityLow = &audio
		} else if audio.Quality == "AUDIO_QUALITY_HIGH" {
			audioQualityHigh = &audio
		}
	}
	if audioQualityHigh != nil {
		return *audioQualityHigh, nil
	} else if audioQualityMedium != nil {
		return *audioQualityMedium, nil
	} else if audioQualityLow != nil {
		return *audioQualityLow, nil
	}
	return Quality{}, errors.New("not found")
}
func (info Streamings) GetHighVideo() (*Quality, error) {
	var video Quality
	toIntVideo := 0
	for _, item := range info.Video {
		toIntVideoNow, _ := strconv.Atoi(item.Quality)
		if toIntVideoNow > toIntVideo {
			video = item
			toIntVideo = toIntVideoNow
		}
	}
	if video.Quality != "" {
		return &video, nil
	}

	return nil, errors.New("not found")
}
func GetVideoInfo(videoId string) (info Info, streamings Streamings, err error) {
	url := "https://m.youtube.com/youtubei/v1/player"

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	requestBody := []byte(fmt.Sprintf(`{"videoId":"%s","context":{"client":{"clientName":"ANDROID_CREATOR","clientVersion":"22.36.102"}}}`, videoId))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	info = Info{}
	streamings = Streamings{}
	if err != nil {
		return info, streamings, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return info, streamings, err
	}
	defer resp.Body.Close()
	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return info, streamings, err
	}

	if data.(map[string]interface{})["playabilityStatus"] != nil && data.(map[string]interface{})["playabilityStatus"].(map[string]interface{})["status"] != nil {
		if data.(map[string]interface{})["playabilityStatus"].(map[string]interface{})["status"].(string) == "LOGIN_REQUIRED" {
			return info, streamings, errors.New("LOGIN_REQUIRED")
		} else if data.(map[string]interface{})["playabilityStatus"].(map[string]interface{})["status"].(string) == "UNPLAYABLE" {
			return info, streamings, errors.New("UNPLAYABLE")
		} else if data.(map[string]interface{})["playabilityStatus"].(map[string]interface{})["status"].(string) == "ERROR" {
			return info, streamings, errors.New("ERROR")
		}
	}
	arrayVideos := data.(map[string]interface{})["streamingData"].(map[string]interface{})["formats"].([]interface{})
	arrayAudios := data.(map[string]interface{})["streamingData"].(map[string]interface{})["adaptiveFormats"].([]interface{})

	var videos []interface{}
	var audios []interface{}
	for _, item := range arrayVideos {
		if strings.Contains(item.(map[string]interface{})["mimeType"].(string), "video/mp4") {
			videos = append(videos, item)
		}
	}
	for _, item := range arrayAudios {
		if strings.Contains(item.(map[string]interface{})["mimeType"].(string), "audio/") {
			audios = append(audios, item)
		}
	}
	var qualitys []Quality
	var qualitysAudio []Quality
	for _, item := range videos {
		qualityLabel := item.(map[string]interface{})["qualityLabel"].(string)
		qualityLabel = strings.Replace(qualityLabel, "p", "", -1)
		mime := item.(map[string]interface{})["mimeType"].(string)
		for _, qualityExists := range qualitys {
			if qualityExists.Quality == qualityLabel {
				continue
			}
		}
		qualitys = append(qualitys, Quality{
			Quality:  qualityLabel,
			Url:      item.(map[string]interface{})["url"].(string),
			MimeType: mime,
		})
	}
	for _, item := range audios {
		qualityLabel := item.(map[string]interface{})["audioQuality"].(string)
		qualityLabel = strings.Replace(qualityLabel, "p", "", -1)
		mime := item.(map[string]interface{})["mimeType"].(string)
		for _, qualityExists := range qualitysAudio {
			if qualityExists.Quality == qualityLabel {
				continue
			}
		}
		qualitysAudio = append(qualitysAudio, Quality{
			Quality:  qualityLabel,
			Url:      item.(map[string]interface{})["url"].(string),
			MimeType: mime,
		})
	}
	details := data.(map[string]interface{})["videoDetails"].(map[string]interface{})
	info = Info{
		Title:       details["title"].(string),
		Thumbnail:   details["thumbnail"].(map[string]interface{})["thumbnails"].([]interface{})[len(details["thumbnail"].(map[string]interface{})["thumbnails"].([]interface{}))-1].(map[string]interface{})["url"].(string),
		ViewCount:   details["viewCount"].(string),
		Author:      details["author"].(string),
		Duration:    details["lengthSeconds"].(string),
	}
	streamings = Streamings{
		Video: qualitys,
		Audio: qualitysAudio,
	}
	sort.Slice(qualitys, func(i, j int) bool {
		return qualitys[i].Quality > qualitys[j].Quality
	})
	return info, streamings, nil
}

var (
	validPathDomains  = regexp.MustCompile(`^https?:\/\/(youtu\.be\/|(www\.)?youtube\.com\/(embed|v|shorts)\/)`)
	validQueryDomains = map[string]bool{
		"www.youtube.com":      true,
		"youtube.com":          true,
		"m.youtube.com":        true,
		"gaming.youtube.com":   true,
		"studio.youtube.com":   true,
		"music.youtube.com":    true,
		"www.youtube-nocookie": true,
	}
	idRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
)

func GetVideoID(link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	id := query.Get("v")
	if validPathDomains.MatchString(link) && id == "" {
		paths := parsed.Path
		if paths != "" {
			// remove leading forward slash if present
			if paths[0] == '/' {
				paths = paths[1:]
			}
			pathParts := regexp.MustCompile(`/`).Split(paths, -1)
			id = pathParts[len(pathParts)-1]
		}
	} else if parsed.Hostname() != "" && !validQueryDomains[parsed.Hostname()] {
		return "", errors.New("not a YouTube domain")
	}
	if id == "" {
		return "", errors.New(`no video id found: "` + link + `"`)
	}
	if len(id) > 11 {
		id = id[:11]
	}
	if !idRegex.MatchString(id) {
		return "", errors.New(`video id (` + id + `) does not match expected format (` + idRegex.String() + `)`)
	}
	return id, nil
}

func ValidateID(id string) bool {
	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
	return idRegex.MatchString(id)
}

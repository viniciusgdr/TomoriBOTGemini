package twitterServices

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tomoribot-geminiai-version/src/utils/hooks"

	"github.com/PuerkitoBio/goquery"
)

type Twitter struct {
	Info struct {
		Img  string
		Text string
	}
	Download []struct {
		URL     string
		Quality string
	}
}

type Result struct {
	Url     string `json:"url"`
	Quality string `json:"quality"`
	ContentType string `json:"contentType"`
}

func getVideoID(url string) string {
	return strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
}

func TwitterOfficial(urlTW string) (Result, error) {
	videoID := getVideoID(urlTW)
	if videoID == "" {
		return Result{}, errors.New("not found")
	}
	request, _ := http.NewRequest("GET", "https://api.twitter.com/graphql/pq4JqttrkAz73WE6s2yUqg/TweetResultByRestId?variables=%7B%22tweetId%22%3A%22"+videoID+"%22%2C%22withCommunity%22%3Afalse%2C%22includePromotedContent%22%3Afalse%2C%22withVoice%22%3Afalse%7D&features=%7B%22creator_subscriptions_tweet_preview_api_enabled%22%3Atrue%2C%22c9s_tweet_anatomy_moderator_badge_enabled%22%3Atrue%2C%22tweetypie_unmention_optimization_enabled%22%3Atrue%2C%22responsive_web_edit_tweet_api_enabled%22%3Atrue%2C%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%22%3Atrue%2C%22view_counts_everywhere_api_enabled%22%3Atrue%2C%22longform_notetweets_consumption_enabled%22%3Atrue%2C%22responsive_web_twitter_article_tweet_consumption_enabled%22%3Atrue%2C%22tweet_awards_web_tipping_enabled%22%3Afalse%2C%22freedom_of_speech_not_reach_fetch_enabled%22%3Atrue%2C%22standardized_nudges_misinfo%22%3Atrue%2C%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%22%3Atrue%2C%22rweb_video_timestamps_enabled%22%3Atrue%2C%22longform_notetweets_rich_text_read_enabled%22%3Atrue%2C%22longform_notetweets_inline_media_enabled%22%3Atrue%2C%22responsive_web_graphql_exclude_directive_enabled%22%3Atrue%2C%22verified_phone_label_enabled%22%3Afalse%2C%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%22%3Afalse%2C%22responsive_web_graphql_timeline_navigation_enabled%22%3Atrue%2C%22responsive_web_enhance_cards_enabled%22%3Afalse%7D&fieldToggles=%7B%22withArticleRichContentState%22%3Atrue%7D", nil)
	headers := map[string]string{
		"accept":                    "*/*",
		"accept-language":           "pt-PT,pt;q=0.9",
		"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
		"content-type":              "application/json",
		"priority":                  "u=1, i",
		"sec-ch-ua":                 "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Linux\"",
		"sec-fetch-dest":            "empty",
		"sec-fetch-mode":            "cors",
		"sec-fetch-site":            "same-site",
		"x-client-transaction-id":   "rTCfTPOcodD+tiDOZuU37mduO5nLSpLa8C47X0320j5CjYOTjbQb7ni/pA98dMyk9cC5I6xzbb6wwRef7bIiRiXTWE/nrA",
		"x-guest-token":             "1762354184156778663",
		"x-twitter-active-user":     "yes",
		"x-twitter-client-language": "pt",
		"cookie":                    "guest_id_marketing=v1%3A170901295509517638; guest_id_ads=v1%3A170901295509517638; guest_id=v1%3A170901295509517638; personalization_id=\"v1_NeHifXW4YXzc45a36y+D4A==\"; gt=1762354184156778663",
		"Referer":                   "https://twitter.com/",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	var resultJSON map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&resultJSON)
	if resultJSON["data"] == nil {
		return Result{}, errors.New("not found")
	}
	tweetResult := resultJSON["data"].(map[string]interface{})["tweetResult"].(map[string]interface{})
	if tweetResult["result"] == nil {
		return Result{}, errors.New("not found")
	}
	result := tweetResult["result"].(map[string]interface{})
	entities := result["legacy"].(map[string]interface{})["entities"].(map[string]interface{})
	media := entities["media"].([]interface{})
	firstMedia := media[0].(map[string]interface{})
	videoInfo := firstMedia["video_info"].(map[string]interface{})
	variants := videoInfo["variants"].([]interface{})
	medias := make([]map[string]interface{}, 0)
	for _, v := range variants {
		media := v.(map[string]interface{})
		if media["content_type"] == "video/mp4" {
			medias = append(medias, media)
		}
	}
	if len(medias) == 0 {
		return Result{}, errors.New("not found")
	}
	// get the highest quality
	var highestQuality int
	var highestQualityIndex int
	for i, v := range medias {
		bitrate := int(v["bitrate"].(float64))
		if bitrate > highestQuality {
			highestQuality = bitrate
			highestQualityIndex = i
		}
	}
	return Result{
		Url:     medias[highestQualityIndex]["url"].(string),
		Quality: strconv.Itoa(highestQuality),
		ContentType: medias[highestQualityIndex]["content_type"].(string),
	}, nil
}

func TwitterGetToken() (tt string, ts string) {
	resp, _ := http.Get("https://ssstwitter.com/")
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	vals, _ := doc.Find(".pure-form.pure-g.hide-after-request").Attr("include-vals")
	vals += "'"
	ttR := hooks.Getstr(vals, "tt:'", "'", 0)
	tsR := strings.Split(hooks.Getstr(vals, "ts:", "'", 0), ",")[0]
	return ttR, tsR
}

func TwitterDownloader(urlTW string) (Result, error) {
	tt, ts := TwitterGetToken()
	param := url.Values{}
	param.Set("id", urlTW)
	param.Set("tt", tt)
	if tt == "" {
		param.Set("tt", "4442124ad689632aafd4a2332e8dc437")
	}
	param.Set("ts", ts)
	if ts == "" {
		param.Set("ts", "1614856639")
	}
	param.Set("locale", "en")
	param.Set("source", "form")

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://ssstwitter.com/", bytes.NewBufferString(param.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.5")
	req.Header.Set("HX-Current-URL", param.Get("id"))
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Target", "target")
	req.Header.Set("Sec-CH-UA", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Brave\";v=\"110\"")
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Cookie", "PHPSESSID=r8t5up1pb50ob3k9b3ghl00tfm; __cflb=0H28v9FyquVAJn5fH3Zv3XFuSyp49MRZwCrAjNaEZwJ")
	req.Header.Set("Referer", param.Get("id"))
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))

	var result Twitter
	result.Info.Img, _ = doc.Find(".result_overlay > img").Attr("src")
	result.Info.Text = doc.Find(".result_overlay > p").Text()
	if result.Info.Text == "🔆 But it contains some image instead, so here you go:" {
		imgSrc, _ := doc.Find("#mainpicture > div > img").Attr("src")
		result.Download = append(result.Download, struct {
			URL     string
			Quality string
		}{URL: imgSrc, Quality: "image"})
	}
	doc.Find(".result_overlay > a").Each(func(i int, s *goquery.Selection) {
		download, _ := s.Attr("href")
		if download == "/" {
			return
		}
		quality := strings.TrimSpace(strings.Replace(s.Text(), "Download ", "", -1))
		result.Download = append(result.Download, struct {
			URL     string
			Quality string
		}{URL: download, Quality: quality})
	})

	var arr []struct {
		URL     string
		Quality string
	}
	for _, v := range result.Download {
		if !strings.Contains(v.URL, "/p/") {
			arr = append(arr, v)
		}
	}

	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			aQuality := strings.Split(strings.Split(arr[i].Quality, "\t")[len(strings.Split(arr[i].Quality, "\t"))-1], "x")
			bQuality := strings.Split(strings.Split(arr[j].Quality, "\t")[len(strings.Split(arr[j].Quality, "\t"))-1], "x")
			aSum, _ := strconv.Atoi(aQuality[0])
			bSum, _ := strconv.Atoi(bQuality[0])
			aSum += bSum
			if len(bQuality) > 1 {
				bSum, _ = strconv.Atoi(bQuality[1])
			}
			aSum -= bSum
			if aSum < 0 {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	if len(arr) != 0 {
		return Result{
			Url:     arr[0].URL,
			Quality: arr[0].Quality,
		}, nil
	} else if len(result.Download) != 0 {
		return Result{
			Url:     result.Download[0].URL,
			Quality: result.Download[0].Quality,
		}, nil
	} else {
		return Result{}, errors.New("not found")
	}
}

package web_functions

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)


func GetJsonFromUrl(url string, headers ...http.Header) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	if (len(headers) > 0) {
		req.Header = headers[0]
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("conteúdo da resposta não é JSON (tipo de conteúdo: %s)", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
func GetBufferFromUrl(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBufferFromUrlThreads(url string) (buffer []byte, sizeFile int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Calculate the size of the file
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	fmt.Println("Size File:", size)
	if err != nil {
		return nil, 0, err
	}

	// Determine the number of downloaders based on file size
	numDownloaders := int(math.Ceil(float64(size) / (2 * 1024 * 1024)))
	if numDownloaders > 16 {
		numDownloaders = 16
	}
	fmt.Println("Num Downloaders:", numDownloaders)
	// Download the file using multiple goroutines
	chunkSize := size / numDownloaders
	chunks := make([][]byte, numDownloaders)

	var wg sync.WaitGroup
	for i := 0; i < numDownloaders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			start := i * chunkSize
			end := (i + 1) * chunkSize
			if i == numDownloaders-1 {
				end = size
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			chunks[i] = body
		}(i)
	}

	// Wait for all the goroutines to complete
	wg.Wait()

	// Concatenate the chunks to get the complete file
	file := make([]byte, size)
	for i := 0; i < numDownloaders; i++ {
		copy(file[i*chunkSize:(i+1)*chunkSize], chunks[i])
	}

	return file, size, nil
}
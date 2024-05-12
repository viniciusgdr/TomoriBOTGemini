package speechtotextServices

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

func SpeechToText(mediaByte []byte) (string, error) {
	url := "https://vmi1001.viniciusgdr.com/speechtotext/api/transcript"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(mediaByte))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/octet-stream")

  res, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer res.Body.Close()
	
  body, err := io.ReadAll(res.Body)
  if err != nil {
    return "", err
  }

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response.Message), nil
}
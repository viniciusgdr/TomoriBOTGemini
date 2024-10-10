package shazamService

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Track struct {
	Layout   string `json:"layout"`
	Type     string `json:"type"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Images   struct {
		Background string `json:"background"`
		Coverart   string `json:"coverart"`
		Coverarthq string `json:"coverarthq"`
		Joecolor   string `json:"joecolor"`
	} `json:"images"`
	Share struct {
		Subject  string `json:"subject"`
		Text     string `json:"text"`
		Href     string `json:"href"`
		Image    string `json:"image"`
		Twitter  string `json:"twitter"`
		Html     string `json:"html"`
		Avatar   string `json:"avatar"`
		Snapchat string `json:"snapchat"`
	} `json:"share"`
	Hub struct {
		Type    string `json:"type"`
		Image   string `json:"image"`
		Actions []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri,omitempty"`
			ID   string `json:"id,omitempty"`
		} `json:"actions"`
		Options []struct {
			Caption string `json:"caption"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri,omitempty"`
			} `json:"actions"`
			Beacondata struct {
				Type         string `json:"type"`
				Providername string `json:"providername"`
			} `json:"beacondata"`
			Image               string `json:"image"`
			Type                string `json:"type"`
			Listcaption         string `json:"listcaption"`
			Overflowimage       string `json:"overflowimage"`
			Colouroverflowimage bool   `json:"colouroverflowimage"`
			Providername        string `json:"providername"`
		} `json:"options"`
		Providers []struct {
			Caption string `json:"caption"`
			Images  struct {
				Overflow string `json:"overflow"`
				Default  string `json:"default"`
			} `json:"images"`
			Actions []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri,omitempty"`
			} `json:"actions"`
			Type string `json:"type"`
		} `json:"providers"`
	} `json:"hub"`
	URL    string `json:"url"`
	Genres struct {
		Primary string `json:"primary"`
	} `json:"genres"`
	Sections interface{} `json:"sections"`
}

type Result struct {
	Matches []struct {
		ID            string  `json:"id"`
		Offset        float64 `json:"offset"`
		Timeskew      float64 `json:"timeskew"`
		Frequencyskew float64 `json:"frequencyskew"`
	} `json:"matches"`
	Location struct {
		Accuracy float64 `json:"accuracy"`
	} `json:"location"`
	Timestamp int64  `json:"timestamp"`
	Timezone  string `json:"timezone"`
	Track     Track  `json:"track"`
	Sections  []struct {
		Type      string `json:"type"`
		Metapages []struct {
			Image   string `json:"image"`
			Caption string `json:"caption"`
		} `json:"metapages"`
		Tabname string `json:"tabname"`
	} `json:"sections"`
}

func ShazamService(path string) (Result, error) {
	var cmd *exec.Cmd
	if (runtime.GOOS == "windows") {
		cmd = exec.Command("python", "./python/recognize.py", path)
	} else {
		cmd = exec.Command("python3", "./python/recognize.py", path)
	}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Erro ao executar o arquivo Python:", err, output)
		os.Exit(1)
	}
	jsonData := string(output)
	var result Result
	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return result, err
	}
	return result, nil
}

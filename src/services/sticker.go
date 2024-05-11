package services

import (
	"encoding/json"
	console "tomoribot-geminiai-version/src/thirdpartyLangs"
)

type StickerResult struct {
	Success   bool  `json:"success"`
	ImagePath string `json:"imagePath"`
}

func AddExifOnSticker(exifPath string, packname string, author string) (StickerResult, error) {
	args := []string{
		"./nodejs/sticker.js",
		"--addexif",
		exifPath,
		packname,
		author,
	}
	result, err := console.RunNodeCommand(args)
	if err != nil {
		return StickerResult{}, err
	}
	var stickerResult interface{}
	err = json.Unmarshal(result.Bytes(), &stickerResult)
	if err != nil {
		return StickerResult{}, err
	}
	var stickerResults StickerResult
	stickerResults.Success = stickerResult.(map[string]interface{})["success"].(bool)
	if !stickerResults.Success {
		return stickerResults, nil
	}
	stickerResults.ImagePath = stickerResult.(map[string]interface{})["imagePath"].(string)
	return stickerResults, nil
}
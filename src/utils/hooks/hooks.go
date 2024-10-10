package hooks

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetPathTemp() string {
	return "./assets/temp/"
}
func GenerateTempFileName(ext string) string {
	ext = strings.Replace(ext, ".", "", -1)
	random := rand.Intn(1000000000)
	pathTemp := GetPathTemp()
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", 0755)
	}
	if _, err := os.Stat(pathTemp); os.IsNotExist(err) {
		os.Mkdir(pathTemp, 0755)
	}
	if _, err := os.Stat(pathTemp + strconv.Itoa(random) + "." + ext); err == nil {
		return GenerateTempFileName(ext)
	}

	return pathTemp + strconv.Itoa(random) + "." + ext
}

func ShuffleString(s string) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune(s)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

func Getstr(s, start, end string, from int) string {
	i := strings.Index(s[from:], start)
	if i == -1 {
		return ""
	}
	i += from
	j := strings.Index(s[i+len(start):], end)
	if j == -1 {
		return ""
	}
	j += i + len(start)
	return s[i+len(start) : j]
}

func BinaryToText(binaryStr string) (string, error) {
	var result string
	for i := 0; i < len(binaryStr); i += 8 {
		if i+8 > len(binaryStr) {
			return "", nil
		}
		bin := binaryStr[i : i+8]
		charCode, err := strconv.ParseInt(bin, 2, 8)
		if err != nil {
			return "", err
		}
		result += string(rune(charCode))
	}

	return result, nil
}

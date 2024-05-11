package services

import (
	"os"
	console "tomoribot-geminiai-version/src/thirdpartyLangs"
)

func RemoveBg(archive string) ([]byte, error) {
	if _, err := os.Stat(archive); os.IsNotExist(err) {
		return nil, err
	}
	archiveResult := archive + ".result.png"
	_, err := console.RunPythonCommand(
		[]string{"./python/removebg.py", archive, archiveResult},
	)
	os.Remove(archive)
	if err != nil {
		return nil, err
	}
	buffer, err := os.ReadFile(archiveResult)
	if err != nil {
		return nil, err
	}
	os.Remove(archiveResult)
	return buffer, nil
}
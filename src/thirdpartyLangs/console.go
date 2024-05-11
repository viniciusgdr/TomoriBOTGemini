package console

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
)

func RunNodeCommand(args []string) (*bytes.Buffer, error) {
	process := exec.Command("node", args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	buf := new(bytes.Buffer)
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		return nil, err
	}

	process.Wait()
	return buf, nil
}

func RunPythonCommand(args []string) (*bytes.Buffer, error) {
	var name string
	if (runtime.GOOS == "windows") {
		name = "python"
	} else {
		name = "python3"
	}
	process := exec.Command(name, args...)
	stdin, err := process.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	buf := new(bytes.Buffer)
	process.Stdout = buf
	process.Stderr = os.Stderr

	if err = process.Start(); err != nil {
		return nil, err
	}

	process.Wait()
	return buf, nil
}
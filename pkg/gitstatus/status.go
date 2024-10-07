package gitstatus

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetDeletedFileContent(path string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("HEAD:%s", path))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func ChangedFiles(path string, ext string) ([]string, []string, error) {
	var m []string
	var d []string

	cmd := exec.Command("git", "status", "-su")
	// set working dir
	cmd.Dir = path
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("exec ( git status -su ) with error: %w\n%s", err, data)
	}
	data = bytes.TrimSpace(data)
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if ext != "" {
			if !bytes.HasSuffix(line, []byte(ext)) {
				continue
			}
		}

		arr := bytes.Split(line, []byte(" "))
		filename := string(arr[len(arr)-1])

		if strings.HasPrefix(filename, "..") {
			continue
		}

		if filename != "" {
			if bytes.HasPrefix(line, []byte("D")) {
				d = append(d, filepath.Join(path, filepath.Join(strings.Split(filename, "/")...)))
			} else {
				m = append(m, filepath.Join(path, filepath.Join(strings.Split(filename, "/")...)))
			}
		}
	}
	return m, d, nil
}

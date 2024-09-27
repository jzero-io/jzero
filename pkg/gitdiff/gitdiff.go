package gitdiff

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetChangedFiles(path string) ([]string, error) {
	return diffFilter(path, "M")
}

func GetDeletedFiles(path string) ([]string, error) {
	return diffFilter(path, "D")
}

func GetAddedFiles(path string) ([]string, error) {
	var files []string
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard", "--directory", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	files = append(files, strings.Split(string(output), "\n")...)
	return files, nil
}

func diffFilter(path, flag string) ([]string, error) {
	var files []string
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("--diff-filter=%s", flag), "HEAD", "--", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	files = append(files, strings.Split(string(output), "\n")...)
	return files, nil
}

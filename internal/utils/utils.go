package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GenerateJwtSecret generates a random JWT secret
func GenerateJwtSecret() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// GenerateAppKey generates a random application key
func GenerateAppKey() (string, error) {
	bytes := make([]byte, 16) // 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// CheckCommandExists checks if a command exists in the system
func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// CompareVersions compares two semantic versions
func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < len(v1Parts) || i < len(v2Parts); i++ {
		var v1Num, v2Num int
		if i < len(v1Parts) {
			v1Num = atoi(v1Parts[i])
		}
		if i < len(v2Parts) {
			v2Num = atoi(v2Parts[i])
		}
		if v1Num > v2Num {
			return 1
		} else if v1Num < v2Num {
			return -1
		}
	}
	return 0
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// ReplaceFile replaces target file with source file
func ReplaceFile(source, target string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(target, input, 0644)
}

// ModifyComposerJSON modifies composer.json to replace ext-swoole with hyperf/engine-swow
func ModifyComposerJSON(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var composer map[string]interface{}
	if err := json.Unmarshal(data, &composer); err != nil {
		return err
	}

	// Modify requires section
	if requires, ok := composer["require"].(map[string]interface{}); ok {
		delete(requires, "ext-swoole")
		requires["hyperf/engine-swow"] = "*"
	}

	// Write back to file
	newData, err := json.MarshalIndent(composer, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, newData, 0644)
}

// GetGitHubFileContent fetches file content from GitHub repository
func GetGitHubFileContent(repo, version, filePath string) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, version, filePath)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

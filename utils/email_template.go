package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadRecoveryEmailTemplate(email, setupPassLink string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	templatePath := filepath.Join(currentDir, "..", "email_templates", "recovery.html")

	htmlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading HTML file: %w", err)
	}

	htmlStr := string(htmlContent)

	htmlStr = strings.ReplaceAll(htmlStr, "[email]", email)
	htmlStr = strings.ReplaceAll(htmlStr, "[setupPassURL]", setupPassLink)

	return htmlStr, nil
}

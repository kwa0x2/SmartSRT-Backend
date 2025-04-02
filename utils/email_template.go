package utils

import (
	"fmt"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadRecoveryEmailTemplate(setupPassLink string) (string, error) {
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

	htmlStr = strings.ReplaceAll(htmlStr, "[setupPassURL]", setupPassLink)

	return htmlStr, nil
}

func LoadContactNotifyTemplate(contact *domain.Contact) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	templatePath := filepath.Join(currentDir, "..", "email_templates", "contact_notify.html")

	htmlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading HTML file: %w", err)
	}

	htmlStr := string(htmlContent)

	htmlStr = strings.ReplaceAll(htmlStr, "[first_name]", contact.FirstName)
	htmlStr = strings.ReplaceAll(htmlStr, "[last_name]", contact.LastName)
	htmlStr = strings.ReplaceAll(htmlStr, "[email]", contact.Email)
	htmlStr = strings.ReplaceAll(htmlStr, "[message]", contact.Message)

	return htmlStr, nil
}

func LoadDeleteAccountEmailTemplate(deleteAccountLink string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	templatePath := filepath.Join(currentDir, "..", "email_templates", "delete_account.html")

	htmlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading HTML file: %w", err)
	}

	htmlStr := string(htmlContent)

	htmlStr = strings.ReplaceAll(htmlStr, "[deleteAccountURL]", deleteAccountLink)

	return htmlStr, nil
}

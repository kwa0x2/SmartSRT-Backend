package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

func loadTemplate(templateName string, replacements map[string]string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	templatePath := filepath.Join(currentDir, "email_templates", templateName)
	htmlContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	htmlStr := string(htmlContent)
	for key, value := range replacements {
		htmlStr = strings.ReplaceAll(htmlStr, key, value)
	}

	return htmlStr, nil
}

func LoadRecoveryEmailTemplate(setupPassLink string) (string, error) {
	return loadTemplate("recovery.html", map[string]string{
		"[setupPassURL]": setupPassLink,
	})
}

func LoadContactNotifyTemplate(contact *domain.Contact) (string, error) {
	return loadTemplate("contact_notify.html", map[string]string{
		"[first_name]": contact.FirstName,
		"[last_name]":  contact.LastName,
		"[email]":      contact.Email,
		"[message]":    contact.Message,
	})
}

func LoadDeleteAccountEmailTemplate(deleteAccountLink string) (string, error) {
	return loadTemplate("delete_account.html", map[string]string{
		"[deleteAccountURL]": deleteAccountLink,
	})
}

func LoadSRTCreatedEmailTemplate(SRTLink string) (string, error) {
	return loadTemplate("srt_notify.html", map[string]string{
		"[SRTLink]": SRTLink,
	})
}

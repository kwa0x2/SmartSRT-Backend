package utils

import "strings"

func ToCamelCase(input string) string {
	if input == "" {
		return input
	}
	return strings.ToUpper(string(input[0])) + strings.ToLower(input[1:])
}

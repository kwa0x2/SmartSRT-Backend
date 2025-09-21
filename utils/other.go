package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToCamelCase(input string) string {
	if input == "" {
		return input
	}
	return strings.ToUpper(string(input[0])) + strings.ToLower(input[1:])
}

// ConvertCentsToDollars converts amount from cents to dollars format (e.g., "499" -> "4.99")
func ConvertCentsToDollars(amountCents string) (string, error) {
	amountInt, err := strconv.Atoi(amountCents)
	if err != nil {
		return "", fmt.Errorf("failed to parse amount: %v", err)
	}
	return fmt.Sprintf("%.2f", float64(amountInt)/100), nil
}

// ParseUserIDFromCustomData extracts and parses user ID from custom data (safe version)
func ParseUserIDFromCustomData(data map[string]interface{}) (bson.ObjectID, error) {
	customData, ok := data["custom_data"].(map[string]interface{})
	if !ok {
		return bson.ObjectID{}, fmt.Errorf("custom_data is not a valid map")
	}

	userIDStr, ok := customData["user_id"].(string)
	if !ok {
		return bson.ObjectID{}, fmt.Errorf("user_id is not a valid string")
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf("invalid user id format: %v", err)
	}

	return userID, nil
}


// ParseProductAndPrice extracts product and price information from items array
func ParseProductAndPrice(data map[string]interface{}) (string, string, string, string, string, error) {
	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", "", "", "", "", fmt.Errorf("invalid items format")
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		return "", "", "", "", "", fmt.Errorf("invalid item format")
	}

	product, ok := item["product"].(map[string]interface{})
	if !ok {
		return "", "", "", "", "", fmt.Errorf("invalid product format")
	}

	price, ok := item["price"].(map[string]interface{})
	if !ok {
		return "", "", "", "", "", fmt.Errorf("invalid price format")
	}

	unitPriceData, ok := price["unit_price"].(map[string]interface{})
	if !ok {
		return "", "", "", "", "", fmt.Errorf("invalid unit_price format")
	}

	amount, err := ConvertCentsToDollars(unitPriceData["amount"].(string))
	if err != nil {
		return "", "", "", "", "", err
	}

	return product["id"].(string), product["name"].(string), price["id"].(string), amount, unitPriceData["currency_code"].(string), nil
}


// ParseBillingPeriod parses billing period from webhook data (returns nil error if not present)
func ParseBillingPeriod(data map[string]interface{}) (time.Time, time.Time, error) {
	billingPeriodData, exists := data["current_billing_period"].(map[string]interface{})
	if !exists {
		return time.Time{}, time.Time{}, nil
	}

	startsAtStr, ok := billingPeriodData["starts_at"].(string)
	if !ok {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid starts_at format in current_billing_period")
	}

	endsAtStr, ok := billingPeriodData["ends_at"].(string)
	if !ok {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid ends_at format in current_billing_period")
	}

	startsAt, err := time.Parse(time.RFC3339, startsAtStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse starts_at: %v", err)
	}

	endsAt, err := time.Parse(time.RFC3339, endsAtStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse ends_at: %v", err)
	}

	return startsAt, endsAt, nil
}

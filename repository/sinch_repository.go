package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"net/http"
)

type SinchRepository struct {
	appKey    string
	appSecret string
}

func NewSinchRepository(appKey, appSecret string) domain.SinchRepository {
	return &SinchRepository{
		appKey:    appKey,
		appSecret: appSecret,
	}
}

func (sr *SinchRepository) SendOTP(phoneNumber string) error {
	url := "https://verification.api.sinch.com/verification/v1/verifications"

	data := map[string]interface{}{
		"identity": map[string]interface{}{
			"type":     "number",
			"endpoint": phoneNumber,
		},
		"method": "sms",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(sr.appKey, sr.appSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en-US")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

func (sr *SinchRepository) VerifyOTP(phoneNumber, code string) (bool, error) {
	url := fmt.Sprintf("https://verification.api.sinch.com/verification/v1/verifications/number/%s", phoneNumber)

	data := map[string]interface{}{
		"method": "sms",
		"sms": map[string]interface{}{
			"code": code,
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(sr.appKey, sr.appSecret)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en-US")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return false, nil
	}

	if resp.StatusCode != 200 {
		return false, errors.New(resp.Status)
	}

	return true, nil
}

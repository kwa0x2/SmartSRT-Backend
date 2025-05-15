package domain

type PaddleWebhookEvent struct {
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
}

type PaddleUseCase interface {
	HandleWebhook(event *PaddleWebhookEvent) error
}

package events

import "github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"

type PaymentCreatedEvent struct {
	EventData internalmodels.Payment `json:"eventdata"`
}

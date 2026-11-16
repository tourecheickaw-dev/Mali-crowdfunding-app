package handlers

import "net/http"

// WebhookVerifier verifies that a webhook indeed comes from the payment aggregator.
// Each aggregator has its own signature method: one implementation per aggregator.
type WebhookVerifier interface {
	Verify(rawBody []byte, headers http.Header) bool
}

// Hub2Verifier: implementation for HUB2 (HUB2-Signature header).
type Hub2Verifier struct {
	Secret string
}

func (v Hub2Verifier) Verify(rawBody []byte, headers http.Header) bool {
	// Without a configured secret, anyone could forge a valid signature: reject everything.
	if v.Secret == "" {
		return false
	}
	return VerifyHub2Signature(rawBody, headers.Get("HUB2-Signature"), v.Secret)
}

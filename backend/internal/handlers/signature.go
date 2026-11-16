package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyHub2Signature verifies the HUB2-Signature header ("s1=...,s0=...") on the raw received body.
// We accept s1 or s0: depending on whether our secret is the new or old one, one of the two will match.
func VerifyHub2Signature(rawBody []byte, header string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(rawBody)
	expected := mac.Sum(nil)

	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 || (kv[0] != "s1" && kv[0] != "s0") {
			continue
		}
		received, err := hex.DecodeString(kv[1])
		if err != nil {
			continue
		}
		if hmac.Equal(expected, received) {
			return true
		}
	}
	return false
}

package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type JWT struct {
	Header    map[string]interface{}
	Payload   map[string]interface{}
	Signature string
}

func ParseToken(token string) (*JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var header map[string]interface{}
	var payload map[string]interface{}

	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}
	return &JWT{
		Header:    header,
		Payload:   payload,
		Signature: parts[2],
	}, nil

}

func (j *JWT) ValidateJWT(secretKey string) error {
	// alg check
	alg, ok := j.Header["alg"].(string)
	if !ok {
		return errors.New("alg not found")
	}
	if alg != "HS256" {
		return errors.New("invalid alg")
	}
	//Re-Encode and payload
	headerJSON, err := json.Marshal(j.Header)
	if err != nil {
		return err
	}
	payloadJSON, err := json.Marshal(j.Payload)
	if err != nil {
		return err
	}
	headerString := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadString := base64.RawURLEncoding.EncodeToString(payloadJSON)

	//signature validation
	data := headerString + "." + payloadString
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if j.Signature != expectedSignature {
		return errors.New("invalid signature")
	}
	// exp check
	if exp, ok := j.Payload["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return errors.New("token expired")
		}
	}
	return nil

}

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type JWT struct {
	Header             *JWTHeader
	Payload            string
	Message, Signature []byte
}

func ParseJWT(input string) (error, *JWT) {
	parts := strings.Split(input, ".")
	decodedParts := make([][]byte, len(parts))
	if len(parts) != 3 {
		return errors.New("invalid jwt: does not contain 3 parts (header, payload, signature)"), nil
	}
	for i := range parts {
		decodedParts[i] = make([]byte, base64.RawURLEncoding.DecodedLen(len(parts[i])))
		if _, err := base64.RawURLEncoding.Decode(decodedParts[i], []byte(parts[i])); err != nil {
			return err, nil
		}
	}
	var parsedHeader JWTHeader
	if err := json.Unmarshal(decodedParts[0], &parsedHeader); err != nil {
		return err, nil
	}
	return nil, &JWT{
		Header:    &parsedHeader,
		Payload:   string(decodedParts[1]),
		Message:   []byte(parts[0] + "." + parts[1]),
		Signature: decodedParts[2],
	}
}

func GenerateSignature(message, secret []byte) []byte {
	hasher := hmac.New(sha256.New, secret)
	hasher.Write(message)
	return hasher.Sum(nil)
}

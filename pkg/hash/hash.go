package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func GenerateHash(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil)), nil
}

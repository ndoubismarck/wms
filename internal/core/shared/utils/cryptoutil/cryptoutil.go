package cryptoutil

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func UID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func UUID() string {
	return uuid.New().String()
}

func Sha1(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	bs := hasher.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func Sha256(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString
}

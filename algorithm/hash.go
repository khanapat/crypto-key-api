package algorithm

import (
	"crypto/sha256"
	"crypto/sha512"

	"github.com/pkg/errors"
)

var (
	errHashType = errors.New("Hash type is not compatible")
)

func hashSha(hashType string, message string) ([]byte, error) {
	switch hashType {
	case "sha224":
		hash := sha256.Sum224([]byte(message))
		return hash[:], nil
	case "sha256":
		hash := sha256.Sum256([]byte(message))
		return hash[:], nil
	case "sha384":
		hash := sha512.Sum384([]byte(message))
		return hash[:], nil
	case "sha512":
		hash := sha512.Sum512([]byte(message))
		return hash[:], nil
	default:
		return nil, errHashType
	}
}

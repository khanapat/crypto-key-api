package algorithm

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func ValidatePublicKeyFromPemStr(pubPEM string) (string, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return "", errors.Wrap(errParseEcdsaPub, parseEcdsaPubTopic)
	}
	publicKey, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
	if err != nil {
		return "", errors.Wrap(err, parseEcdsaPubTopic)
	}
	switch publicKey.(type) {
	case *ecdsa.PublicKey:
		return "ECDSA", nil
	case *rsa.PublicKey:
		return "RSA", nil
	default:
		return "undefined type", errors.Wrap(errCastTypeEcdsaPub, parseEcdsaPubTopic)
	}
}

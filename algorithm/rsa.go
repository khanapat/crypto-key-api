package algorithm

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
)

const (
	signRsaTopic     = "Signing Rsa"
	verifyRsaTopic   = "Verifying Rsa"
	parseRsaPriTopic = "Parsing Private Key Rsa"
	parseRsaPubTopic = "Parsing Public Key Rsa"
)

var (
	errParseRsaPri    = errors.New("Failed to parse RSA private key")
	errParseRsaPub    = errors.New("Failed to parse RSA public key")
	errCastTypeRsaPub = errors.New("Unsupported public key type")
)

func SignRsa(ctx context.Context, privateKey *rsa.PrivateKey, hashType, message string) ([]byte, *string, error) {
	logger := middleware.ContextData(ctx)

	hash, err := hashSha(hashType, message)
	if err != nil {
		return nil, nil, errors.Wrap(err, hashTypeTopic)
	}
	logger.Debug(fmt.Sprintf("Hash - Base64 %q Hex \"%x\"\n", bToB64(hash), hash))
	sigHex, err := rsa.SignPKCS1v15(rand.Reader, privateKey, cryptoHash(hashType), hash)
	if err != nil {
		return nil, nil, errors.Wrap(err, signRsaTopic)
	}
	sigB64 := bToB64(sigHex)
	logger.Debug(fmt.Sprintf("Signature - Base64 %q Hex \"%x\"\n", sigB64, sigHex))
	return sigHex, &sigB64, nil
}

func VerifyRsa(ctx context.Context, publicKey *rsa.PublicKey, hashType, message, sig string) (bool, error) {
	logger := middleware.ContextData(ctx)

	hash, err := hashSha(hashType, message)
	if err != nil {
		return false, errors.Wrap(err, hashTypeTopic)
	}
	logger.Debug(fmt.Sprintf("Hash - Base64 %q Hex \"%x\"\n", bToB64(hash), hash))
	sigB64, _ := b64ToB(sig)
	err = rsa.VerifyPKCS1v15(publicKey, cryptoHash(hashType), hash, sigB64)
	if err != nil {
		return false, errors.Wrap(err, verifyRsaTopic)
	}
	return true, nil
}

func ParseRsaPrivateKeyFromPemStr(priPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priPEM))
	if block == nil {
		return nil, errors.Wrap(errParseRsaPri, parseRsaPriTopic)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(block.Bytes))
	if err != nil {
		return nil, errors.Wrap(err, parseRsaPriTopic)
	}
	return privateKey, nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.Wrap(errParseRsaPub, parseRsaPubTopic)
	}
	publicKey, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
	if err != nil {
		return nil, errors.Wrap(err, parseRsaPubTopic)
	}
	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.Wrap(errCastTypeRsaPub, parseRsaPubTopic)
	}
}

func cryptoHash(hashType string) crypto.Hash {
	var cryp crypto.Hash
	switch hashType {
	case "sha224":
		cryp = crypto.SHA224
	case "sha256":
		cryp = crypto.SHA256
	case "sha384":
		cryp = crypto.SHA384
	case "sha512":
		cryp = crypto.SHA512
	}
	return cryp
}

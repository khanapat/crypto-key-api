package algorithm

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
)

const (
	hashTypeTopic        = "Hashing"
	generateKeyTopic     = "Generating Private & Public Key"
	signEcdsaTopic       = "Signing ECDSA"
	parseEcdsaPriTopic   = "Parsing Private Key ECDSA"
	parseEcdsaPubTopic   = "Parsing Public Key ECDSA"
	marshalEcdsaPriTopic = "Marshalling Private Key ECDSA"
	marshalEcdsaPubTopic = "Marshalling Public Key ECDSA"
)

var (
	errParseEcdsaPri    = errors.New("Failed to parse ECDSA private key")
	errParseEcdsaPub    = errors.New("Failed to parse ECDSA public key")
	errCastTypeEcdsaPub = errors.New("Unsupported public key type")
)

func GenerateEcdsaKey(ctx context.Context, curveType string) (*ecdsa.PrivateKey, error) {
	logger := middleware.ContextData(ctx)

	privateKey, err := ecdsa.GenerateKey(curveHash(curveType), rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, generateKeyTopic)
	}
	logger.Debug(fmt.Sprintf("PrivateKey - %v", privateKey))
	return privateKey, nil
}

func SignEcdsa(ctx context.Context, privateKey *ecdsa.PrivateKey, hashType, message string) ([]byte, *string, error) {
	logger := middleware.ContextData(ctx)

	hash, err := hashSha(hashType, message)
	if err != nil {
		return nil, nil, errors.Wrap(err, hashTypeTopic)
	}
	logger.Debug(fmt.Sprintf("Hash - Base64 %q Hex \"%x\"\n", bToB64(hash), hash))
	sigHex, err := ecdsa.SignASN1(rand.Reader, privateKey, hash)
	if err != nil {
		return nil, nil, errors.Wrap(err, signEcdsaTopic)
	}
	sigB64 := bToB64(sigHex)
	logger.Debug(fmt.Sprintf("Signature - Base64 %q Hex \"%x\"\n", sigB64, sigHex))
	return sigHex, &sigB64, nil
}

func VerifyEcdsa(ctx context.Context, publicKey *ecdsa.PublicKey, hashType, message, sig string) (bool, error) {
	logger := middleware.ContextData(ctx)

	hash, err := hashSha(hashType, message)
	if err != nil {
		return false, errors.Wrap(err, hashTypeTopic)
	}
	logger.Debug(fmt.Sprintf("Hash - Base64 %q Hex \"%x\"\n", bToB64(hash), hash))
	sigB64, _ := b64ToB(sig)
	valid := ecdsa.VerifyASN1(publicKey, hash, sigB64)
	logger.Debug(fmt.Sprintf("Signature verified - %t\n", valid))
	return valid, nil
}

func ParseEcdsaPrivateKeyFromPemStr(priPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priPEM))
	if block == nil {
		return nil, errors.Wrap(errParseEcdsaPri, parseEcdsaPriTopic)
	}
	privateKey, err := x509.ParseECPrivateKey([]byte(block.Bytes))
	if err != nil {
		return nil, errors.Wrap(err, parseEcdsaPriTopic)
	}
	return privateKey, nil
}

func ParseEcdsaPublicKeyFromPemStr(pubPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.Wrap(errParseEcdsaPub, parseEcdsaPubTopic)
	}
	publicKey, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
	if err != nil {
		return nil, errors.Wrap(err, parseEcdsaPubTopic)
	}
	switch pub := publicKey.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.Wrap(errCastTypeEcdsaPub, parseEcdsaPubTopic)
	}
}

func MarshalEcdsaPrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	privByte, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, marshalEcdsaPriTopic)
	}
	pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privByte,
		},
	)
	return pem, nil
}

func MarshalEcdsaPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	pubByte, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, marshalEcdsaPubTopic)
	}
	pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PUBLIC KEY",
			Bytes: pubByte,
		},
	)
	return pem, nil
}

func curveHash(hashType string) elliptic.Curve {
	var curve elliptic.Curve
	switch hashType {
	case "p224":
		curve = elliptic.P224()
	case "p256":
		curve = elliptic.P256()
	case "p384":
		curve = elliptic.P384()
	case "p512":
		curve = elliptic.P521()
	}
	return curve
}

func bToB64(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

func b64ToB(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

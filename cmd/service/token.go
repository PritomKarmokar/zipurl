package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClient struct {
	jwtVerifyingKey    *rsa.PublicKey
	jwtSigningKey      *rsa.PrivateKey
	jwtAlgorithm       string
	jwtAudience        []string
	jwtIssuer          string
	accessTokenExpiry  int64
	refreshTokenExpiry int64
}

type JWTException struct {
	Message string
}

type CustomClaims struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
	jwt.RegisteredClaims
}

func (e *JWTException) Error() string { return e.Message }

func NewJWTException(message string) *JWTException {
	return &JWTException{Message: message}
}

func NewTokenClient(
	jwtVerifyingKey string,
	jwtSigningKey string,
	jwtAlgorithm string,
	jwtAudience []string,
	jwtIssuer string,
	accessTokenExpiry int64,
	refreshTokenExpiry int64,
) (*TokenClient, error) {
	// Load public key for verification
	publicKey, err := loadPublicKeyFromBase64(jwtVerifyingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load verifying key: %w", err)
	}

	// Load private key for signing
	privateKey, err := loadPrivateKeyFromBase64(jwtSigningKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load signing key: %w", err)
	}
	return &TokenClient{
		jwtVerifyingKey:    publicKey,
		jwtSigningKey:      privateKey,
		jwtAlgorithm:       jwtAlgorithm,
		jwtAudience:        jwtAudience,
		jwtIssuer:          jwtIssuer,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}, nil
}

// loadPublicKeyFromBase64 loads an RSA public key from a base64-encoded PEM string
func loadPublicKeyFromBase64(keyB64 string) (*rsa.PublicKey, error) {
	// Decode base64
	keyPEM, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Parse Public Key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaKey, nil
}

func loadPrivateKeyFromBase64(keyB64 string) (*rsa.PrivateKey, error) {
	// Decode base64
	keyPEM, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Parse private key
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS1 format if PKCS8 fails
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return rsaKey, nil
}

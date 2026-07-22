package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

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

// getSigningMethod returns the JWT signing method based on the algorithm string
func (tc *TokenClient) getSigningMethod() jwt.SigningMethod {
	switch tc.jwtAlgorithm {
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS384":
		return jwt.SigningMethodRS384
	case "RS512":
		return jwt.SigningMethodRS512
	default:
		return jwt.SigningMethodRS512
	}
}

// CreateToken creates a JWT token with the given payload and token type
func (tc *TokenClient) CreateToken(payload map[string]interface{}, tokenType string, expirySeconds int64) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		Type: tokenType,
		Data: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateJTI(), // Generate unique token ID for blacklisting
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expirySeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    tc.jwtIssuer,
			Audience:  tc.jwtAudience,
		},
	}
	token := jwt.NewWithClaims(tc.getSigningMethod(), claims)

	signedToken, err := token.SignedString(tc.jwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

// CreateAccessAndRefreshTokens generates both access and refresh tokens
func (tc *TokenClient) CreateAccessAndRefreshTokens(payload map[string]interface{}) (*TokenPair, error) {
	// Create a copy of payload for access token
	accessPayload := make(map[string]interface{})
	for k, v := range payload {
		accessPayload[k] = v
	}
	accessToken, err := tc.CreateToken(accessPayload, "access", tc.accessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create a copy of payload for refresh token
	refreshPayload := make(map[string]interface{})
	for k, v := range payload {
		refreshPayload[k] = v
	}

	refreshToken, err := tc.CreateToken(refreshPayload, "refresh", tc.refreshTokenExpiry)
	return &TokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

// RefreshAccessToken validates a refresh token and returns a new access token
func (tc *TokenClient) RefreshAccessToken(refreshToken string) (map[string]interface{}, error) {
	// Verify the refresh token
	return nil, nil
}

// VerifyToken verifies and decodes a JWT token
func (tc *TokenClient) VerifyToken(tokenString string, tokenType string) (map[string]interface{}, error) {
	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tc.jwtVerifyingKey, nil
	})

	if err != nil {
		return nil, NewJWTException(fmt.Sprintf("Token verification failed: %v", err))
	}

	if !token.Valid {
		return nil, NewJWTException("Invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, NewJWTException("Invalid token claims")
	}

	// Verify token type
	if claims.Type != tokenType {
		return nil, NewJWTException(fmt.Sprintf("Invalid token type, expected %s", tokenType))
	}

	// Verify audience
	audienceMatched := false
	for _, aud := range claims.Audience {
		for _, expectedAud := range tc.jwtAudience {
			if aud == expectedAud {
				audienceMatched = true
				break
			}
		}
		if audienceMatched {
			break
		}
	}

	if !audienceMatched {
		return nil, NewJWTException("Invalid audience")
	}

	// Verify Issuer
	if claims.Issuer != tc.jwtIssuer {
		return nil, NewJWTException("Invalid issuer")
	}

	// Convert claims to map
	result := map[string]interface{}{
		"type": claims.Type,
		"data": claims.Data,
		"jti":  claims.ID,
		"exp":  claims.ExpiresAt,
		"iat":  claims.IssuedAt,
		"iss":  claims.Issuer,
		"aud":  claims.Audience,
	}
	return result, nil
}

// generateJTI generates a unique JWT ID for token identification and blacklisting
func generateJTI() string {
	b := make([]byte, 16) // 128-bit random ID
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

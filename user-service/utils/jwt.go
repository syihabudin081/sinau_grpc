package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// JWTManager holds configuration for JWT token generation
type JWTManager struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// NewJWTManager creates a new token configuration
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		SecretKey:       secretKey,
		AccessTokenTTL:  24 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
}

// GenerateAccessToken creates a new access token
func (jw *JWTManager) GenerateAccessToken(userID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jw.AccessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jw.SecretKey))
}

// GenerateRefreshToken creates a new refresh token
func (jw *JWTManager) GenerateRefreshToken() (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(jw.RefreshTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jw.SecretKey))
}

// ValidateToken checks if the token is valid and returns its claims
func (jw *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jw.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Extract and type assert claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshAccessToken generates a new access token using refresh token claims
func (jw *JWTManager) RefreshAccessToken(refreshToken string, userID, role string) (string, error) {
	// First validate the refresh token
	_, err := jw.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Generate a new access token
	return jw.GenerateAccessToken(userID, role)
}

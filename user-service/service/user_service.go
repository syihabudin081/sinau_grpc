package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"user-service/db"
	"user-service/pb"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	db *db.AuthProvider
	pb.UnimplementedAuthServiceServer
}

const (
	tokenExpiration    = 24 * time.Hour
	refreshTokenExpiry = 7 * 24 * time.Hour
	jwtSigningKey      = "your-secret-key" // In production, use environment variable
)

func NewAuthServiceHandler(authDB *db.AuthProvider) *AuthServiceHandler {
	return &AuthServiceHandler{db: authDB}
}

func (s *AuthServiceHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Password hashing failed")
	}

	now := time.Now()
	user := &pb.UserORM{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	// Call database method to create user
	_, err = s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *AuthServiceHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Retrieve user by username
	user, err := s.db.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	// Generate tokens
	token, refreshToken, err := generateTokens(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Token generation failed")
	}

	return &pb.LoginResponse{
		UserId:       fmt.Sprintf("%d", user.Id),
		Token:        token,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (s *AuthServiceHandler) Validate(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := validateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		IsValid: true,
		UserId:  claims.UserId,
		Role:    claims.Role,
	}, nil
}

// JWT Claims
type JWTClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func (s *AuthServiceHandler) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	var token string

	// Jika req kosong, ambil token dari metadata
	if req == nil {
		// Ambil metadata dari konteks
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "No metadata found in request")
		}
		log.Printf("Metadata: %v", md)

		// Ambil token dari header "authorization"
		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token is missing")
		}

		// Ambil token dari header
		token = authHeader[0]
	} else {
		// Jika req tidak kosong, kamu bisa mengakses data dari req
		// Misalnya, kamu bisa menggunakan req.UserID atau lainnya jika perlu
		token = req.Token // Pastikan `req` berisi token jika diperlukan
	}

	// Validasi token
	log.Println(token)
	claims, err := validateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	// Ambil data user berdasarkan UserID dari klaim token
	user, err := s.db.FindUserByID(ctx, claims.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Return data user
	return &pb.GetMeResponse{
		IsValid:   true,
		IsExpired: false,
		UserID:    user.Id,
		Username:  user.Username,
		Role:      user.Role,
	}, nil
}

func generateTokens(user *pb.UserORM) (string, string, error) {
	// Generate Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserId: fmt.Sprintf("%d", user.Id),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiration).Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenExpiry).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	authPB "transaction-service/lib/stubs/user_pb"
)

type AuthInterceptor struct {
	authClient authPB.AuthServiceClient
}

// NewAuthInterceptor membuat instance baru dari AuthInterceptor
func NewAuthInterceptor(authClient authPB.AuthServiceClient) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient}
}

// UnaryInterceptor memvalidasi setiap permintaan unary
func (a *AuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Ambil metadata dari context
	log.Print(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	//

	//log.Print("MD : ", md)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	// Ambil token dari metadata
	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
	}
	token := authHeader[0]
	// Hapus awalan "Bearer " jika ada
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	// Validasi token dengan AuthService
	resp, err := a.authClient.GetMe(ctx, &authPB.GetMeRequest{Token: token})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Simpan informasi user dalam context
	ctx = context.WithValue(ctx, "userRole", resp.Role)
	ctx = context.WithValue(ctx, "userID", resp.UserID)

	// Lanjutkan ke handler asli
	return handler(ctx, req)
}

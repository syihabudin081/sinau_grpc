package service

//
//import (
//	"context"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/metadata"
//	"google.golang.org/grpc/status"
//	"log"
//	jwt "user-service/utils"
//)
//
//type AuthInterceptor struct {
//	jwtManager       *jwt.JWTManager
//	accesibleMethods map[string][]string
//}
//
//func NewAuthInterceptor(jwtManager *jwt.JWTManager, accesibleMethods map[string][]string) *AuthInterceptor {
//	return &AuthInterceptor{jwtManager, accesibleMethods}
//}
//
//func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
//	return func(
//		ctx context.Context,
//		req interface{},
//		info *grpc.UnaryServerInfo,
//		handler grpc.UnaryHandler,
//	) (interface{}, error) {
//		log.Println("--> unary interceptor: ", info.FullMethod)
//
//		// TODO: implement authorization
//
//		return handler(ctx, req)
//	}
//}
//
//func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
//	return func(
//		srv interface{},
//		stream grpc.ServerStream,
//		info *grpc.StreamServerInfo,
//		handler grpc.StreamHandler,
//	) error {
//		log.Println("--> stream interceptor: ", info.FullMethod)
//
//		// TODO: implement authorization
//
//		return handler(srv, stream)
//	}
//}
//
//func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
//	accessibleRoles, ok := interceptor.accesibleMethods[method]
//	if !ok {
//		return nil
//	}
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
//	}
//
//	values := md["authorization"]
//	if len(values) == 0 {
//		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
//	}
//
//	accessToken := values[0]
//	claims, err := interceptor.jwtManager.ValidateToken(accessToken)
//	if err != nil {
//		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
//	}
//
//	for _, role := range accessibleRoles {
//		if role == claims.Role {
//			return nil
//		}
//	}
//
//	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
//
//}
//
//func AccessibleMethods() map[string][]string {
//	// TODO : only admin can access the addProduct method
//	return map[string][]string{
//		"/pb.ProductService/AddProduct":    {"admin"},
//		"/pb.ProductService/GetProduct":    {"user", "admin"},
//		"/pb.ProductService/UpdateProduct": {"admin"},
//		"/pb.ProductService/DeleteProduct": {"admin"},
//	}
//}

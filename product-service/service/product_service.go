package service

import (
	"context"
	"go-grpc-gateway/db"
	authPB "go-grpc-gateway/lib/stubs/user_pb"
	"go-grpc-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// Embed the UnimplementedProductServiceServer to satisfy the interface.
type ProductServiceHandler struct {
	db                                   *db.GormProvider
	pb.UnimplementedProductServiceServer // Embed this struct
	authClient                           authPB.AuthServiceClient
}

func NewProductServiceHandler(db *db.GormProvider, authConn *grpc.ClientConn) *ProductServiceHandler {
	return &ProductServiceHandler{
		db:         db,
		authClient: authPB.NewAuthServiceClient(authConn),
	}
}

func (s *ProductServiceHandler) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {

	// log the context
	log.Println("Context: ", ctx)

	// get the metadata from the context

	// check the userRole from the metadata
	userRole, ok := ctx.Value("userRole").(string)

	// check if the userRole is not found
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing user role")
	}
	// 

	// check if the userRole is admin
	if userRole != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "only admin can access this RPC")
	}

	resp, err := s.db.AddProduct(ctx, req)
	if err != nil {
		// Handle the error and return a grpc error or custom error
		return nil, err
	}
	return resp, nil
}

func (s *ProductServiceHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	// Call the GetProduct method from GormProvider
	resp, err := s.db.GetProduct(ctx, req)
	if err != nil {
		// Handle the error and return a grpc error or custom error
		return nil, err
	}
	return resp, nil
}

func (s *ProductServiceHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	// Call the UpdateProduct method from GormProvider
	resp, err := s.db.UpdateProduct(ctx, req)
	if err != nil {
		// Handle the error and return a grpc error or custom error
		return nil, err
	}
	return resp, nil
}

func (s *ProductServiceHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	// Call the DeleteProduct method from GormProvider
	resp, err := s.db.DeleteProduct(ctx, req)
	if err != nil {
		// Handle the error and return a grpc error or custom error
		return nil, err
	}
	return resp, nil
}

func (s *ProductServiceHandler) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	// Call the ListProduct method from GormProvider
	resp, err := s.db.ListProduct(ctx, req)
	if err != nil {
		// Handle the error and return a grpc error or custom error
		return nil, err
	}
	return resp, nil
}

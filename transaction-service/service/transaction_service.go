package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"transaction-service/db"
	productPB "transaction-service/lib/stubs/product_pb"
	authPB "transaction-service/lib/stubs/user_pb"
	"transaction-service/pb"
)

type TransactionServiceHandler struct {
	db *db.TransactionProvider
	pb.UnimplementedTransactionServiceServer
	authClient    authPB.AuthServiceClient
	productClient productPB.ProductServiceClient
}

func NewTransactionServiceHandler(db *db.TransactionProvider, authConn *grpc.ClientConn, productConn *grpc.ClientConn) *TransactionServiceHandler {
	return &TransactionServiceHandler{
		db:            db,
		authClient:    authPB.NewAuthServiceClient(authConn),
		productClient: productPB.NewProductServiceClient(productConn),
	}
}

func (s *TransactionServiceHandler) AddTransaction(ctx context.Context, req *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {

	productIdStr := strconv.FormatUint(req.ProductId, 10)
	md, _ := metadata.FromIncomingContext(ctx)
	newCtx := metadata.NewOutgoingContext(context.Background(), md)
	log.Print("New Context : ", newCtx)
	res, err := s.productClient.GetProduct(newCtx, &productPB.GetProductRequest{ProductId: productIdStr})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product id not found: %v", err)
	}
	log.Print("Product With ID : %v Found With Response : %v", productIdStr, res)
	resp, err := s.db.CreateTransaction(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add transaction: %v", err)
	}
	return resp, nil
}

func (s *TransactionServiceHandler) ApproveTransaction(ctx context.Context, req *pb.ApproveTransactionRequest) (*pb.ApproveTransactionResponse, error) {
	userRole, ok := ctx.Value("userRole").(string)
	// check if the userRole is not found
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing user role")
	}

	// check if the userRole is admin
	if userRole != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "only admin can access this....")
	}
	resp, err := s.db.ApproveTransaction(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve transaction: %v", err)
	}
	return resp, nil
}

func (s *TransactionServiceHandler) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {

	resp, err := s.db.GetTransaction(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transaction: %v", err)
	}
	return resp, nil
}

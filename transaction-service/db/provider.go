package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
	"transaction-service/pb"
)

type TransactionProvider struct {
	db *gorm.DB
}

func NewTransactionProvider(db *gorm.DB) *TransactionProvider {
	return &TransactionProvider{db: db}
}

// GetTransaction fetches a transaction by its ID.
func (p *TransactionProvider) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {

	var transaction pb.TransactionORM
	if err := p.db.WithContext(ctx).First(&transaction, "id = ?", req.TransactionId).Error; err != nil {
		return nil, err
	}

	return &pb.GetTransactionResponse{
		ProductId: transaction.ProductId,
		UserId:    transaction.UserId,
		Quantity:  transaction.Quantity,
		Status:    transaction.Status,
		Note:      transaction.Note,
		CreatedAt: transaction.CreatedAt.String(),
		UpdatedAt: transaction.UpdatedAt.String(),
	}, nil

}

func (p *TransactionProvider) CreateTransaction(ctx context.Context, req *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {
	// Validasi nilai status apakah sesuai dengan enum
	if req.Status < pb.Status_STATUS_UNSPECIFIED || req.Status > pb.Status_STATUS_CANCELLED {
		return nil, fmt.Errorf("invalid status value")
	}
	userID := ctx.Value("userID").(uint64)
	// call product service

	// fimd

	now := time.Now()
	transaction := &pb.TransactionORM{
		ProductId: req.ProductId,
		UserId:    userID,
		Quantity:  req.Quantity,
		Status:    1,
		Note:      req.Note,
		CreatedAt: &now,
		UpdatedAt: nil,
	}

	if err := p.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return nil, err
	}
	return &pb.AddTransactionResponse{
		Message: "Transaction added successfully",
	}, nil
}

func (p *TransactionProvider) ApproveTransaction(ctx context.Context, req *pb.ApproveTransactionRequest) (*pb.ApproveTransactionResponse, error) {

	// TODO : add query for approve a transaction
	if err := p.db.WithContext(ctx).Model(&pb.TransactionORM{}).Debug().Where("id = ?", req.TransactionId).Update("status", 2).Error; err != nil {
		return nil, err
	}

	return &pb.ApproveTransactionResponse{
		Message: fmt.Sprintf("Transaction With ID: %v Approved", req.GetTransactionId()),
	}, nil
}

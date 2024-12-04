package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"transaction-service/pb"
)

type TransactionProvider struct {
	db *sql.DB
}

func NewTransactionProvider(db *sql.DB) *TransactionProvider {
	return &TransactionProvider{db: db}
}

// GetTransaction fetches a transaction by its ID.
func (p *TransactionProvider) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	// Buat query builder dengan query dasar
	qb := NewQueryBuilderNative("SELECT product_id, user_id, quantity, status, note, created_at, updated_at FROM transactions")

	// Tambahkan filter secara dinamis
	qb.Scope("id = ?", req.TransactionId).
		Scope("status = ?", req.Status).
		Scope("user_id = ?", ctx.Value("userID")).
		OrderBy("created_at", "DESC").
		Pagination(1, 0) // Ambil satu data

	query, args := qb.Build()

	// Eksekusi query
	row := p.db.QueryRowContext(ctx, query, args...)
	var transaction pb.TransactionORM
	var updatedAt sql.NullTime
	if err := row.Scan(
		&transaction.ProductId,
		&transaction.UserId,
		&transaction.Quantity,
		&transaction.Status,
		&transaction.Note,
		&transaction.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}

	// Konversi hasil
	transaction.UpdatedAt = nil
	if updatedAt.Valid {
		transaction.UpdatedAt = &updatedAt.Time
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

	// Ambil userID dari context
	userID := ctx.Value("userID").(uint64)

	// Waktu sekarang
	now := time.Now()

	// Membangun query menggunakan QueryBuilderNative
	qb := NewQueryBuilderNative("INSERT INTO transactions (product_id, user_id, quantity, status, note, created_at, updated_at)").
		Scope("?", req.ProductId).
		Scope("?", userID).
		Scope("?", req.Quantity).
		Scope("?", req.Status).
		Scope("?", req.Note).
		Scope("?", now).
		Scope("?", nil) // updated_at bisa NULL

	// Membangun query dan args
	query, args := qb.Build()

	// Eksekusi query
	_, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return &pb.AddTransactionResponse{
		Message: "Transaction added successfully",
	}, nil
}

func (p *TransactionProvider) ApproveTransaction(ctx context.Context, req *pb.ApproveTransactionRequest) (*pb.ApproveTransactionResponse, error) {
	// Buat query builder
	qb := NewQueryBuilderNative("UPDATE transactions SET status = ?").
		Scope("id = ?", req.TransactionId)

	query, args := qb.Build()

	// Tambahkan nilai status
	args = append([]interface{}{2}, args...) // Status = 2 (Approved)

	// Eksekusi query
	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return nil, err
	}

	return &pb.ApproveTransactionResponse{
		Message: fmt.Sprintf("Transaction With ID: %v Approved", req.GetTransactionId()),
	}, nil
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
func (p *TransactionProvider) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.ListTransactionResponse, error) {
	// Start building the query.
	query := `
        SELECT id, product_id, user_id, quantity, status, note, created_at, updated_at 
        FROM transaction_orms 
        WHERE 1=1 
    `
	var args []interface{}
	argCount := 1

	// Add filters based on the request parameters.
	if req.GetTransactionId() != "" {
		query += fmt.Sprintf(" AND id = $%d", argCount)
		args = append(args, req.GetTransactionId())
		argCount++
	}
	if req.GetStatus() != -1 {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, req.GetStatus())
		argCount++
	}

	// Order by creation date.
	query += ` ORDER BY created_at DESC;`

	// Execute the query to fetch the transactions.
	rows, err := p.db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Collect the results.
	var transactions []*pb.GetTransactionResponse
	for rows.Next() {
		var transaction pb.GetTransactionResponse
		var status int32
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.ProductId,
			&transaction.UserId,
			&transaction.Quantity,
			&status,
			&transaction.Note,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Convert integer status to the enum.
		transaction.Status = pb.Status(status)

		// Add the transaction to the response list.
		transactions = append(transactions, &transaction)
	}

	// Check for errors in row iteration.
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// Return the list of transactions.
	return &pb.ListTransactionResponse{
		Transactions: transactions,
	}, nil
}

func (p *TransactionProvider) CreateTransaction(ctx context.Context, transaction *pb.AddTransactionRequest) (*pb.AddTransactionResponse, error) {

	// get user id from context
	user_id := ctx.Value("userID").(uint64)

	query := `
		INSERT INTO transaction_orms (product_id, user_id, quantity, status, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	var transactionID int64
	err := p.db.QueryRow(query, transaction.ProductId, user_id, transaction.Quantity, transaction.Status, transaction.Note, time.Now(), time.Now()).Scan(&transactionID)
	if err != nil {
		log.Printf("Error adding transaction: %v", err)
		return nil, err
	}
	log.Printf("Transaction added with ID: %d", transactionID)
	return &pb.AddTransactionResponse{
		Message: "Transaction added successfully",
	}, nil
}

func (p *TransactionProvider) ApproveTransaction(ctx context.Context, req *pb.ApproveTransactionRequest) (*pb.ApproveTransactionResponse, error) {
	// Define the SQL query to update the transaction status.
	query := `UPDATE transaction_orms SET status = $1 WHERE id = $2;`

	// Execute the query.
	result, err := p.db.ExecContext(ctx, query, pb.Status_STATUS_COMPLETED, req.GetTransactionId())
	if err != nil {
		log.Printf("Error approving transaction: %v", err)
		return nil, fmt.Errorf("failed to approve transaction: %w", err)
	}

	// Check if any rows were affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return nil, fmt.Errorf("failed to verify update: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("transaction with ID: %v not found or already updated", req.GetTransactionId())
	}

	// Return a success message.
	return &pb.ApproveTransactionResponse{
		Message: fmt.Sprintf("Transaction with ID: %v approved successfully", req.GetTransactionId()),
	}, nil
}

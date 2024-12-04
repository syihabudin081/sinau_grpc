package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"transaction-service/db"
	authPB "transaction-service/lib/stubs/user_pb"
	"transaction-service/pb"
	"transaction-service/service"
)

func main() {
	postgresDB := db.ConnectDatabase()
	//if err := postgresDB.AutoMigrate(&pb.TransactionORM{}); err != nil {
	//	log.Fatalf("Failed to migrate database: %v", err)
	//}

	// Initialize DB Provider
	dbProvider := db.NewTransactionProvider(postgresDB)

	// gRPC
	authConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	productConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer productConn.Close()

	// Buat AuthInterceptor
	authClient := authPB.NewAuthServiceClient(authConn)
	authInterceptor := service.NewAuthInterceptor(authClient)

	// gRPC Server
	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor),
	)

	// Register TransactionService
	transactionService := service.NewTransactionServiceHandler(dbProvider, authConn, productConn)
	pb.RegisterTransactionServiceServer(grpcServer, transactionService)

	// Run gRPC server in a goroutine
	go func() {
		log.Println("ProductService running on port 50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// gRPC-Gateway Server
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterTransactionServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50054", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway handler: %v", err)
	}

	// Start HTTP server for the gRPC-Gateway
	log.Println("gRPC-Gateway server running on port 8083")
	if err := http.ListenAndServe(":8083", gwmux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}

}

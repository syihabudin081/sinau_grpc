package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-grpc-gateway/db"
	authPB "go-grpc-gateway/lib/stubs/user_pb"
	"go-grpc-gateway/pb"
	"go-grpc-gateway/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
)

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=123 dbname=product_dev_db port=5432 sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate models
	if err := gormDB.AutoMigrate(&pb.ProductORM{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize GormProvider
	dbProvider := db.NewGormProvider(gormDB)

	// gRPC connection to AuthService
	authConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()
	// Buat AuthInterceptor
	authClient := authPB.NewAuthServiceClient(authConn)
	authInterceptor := service.NewAuthInterceptor(authClient)

	// gRPC Server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor),
	)

	// Register ProductService
	productService := service.NewProductServiceHandler(dbProvider, authConn)
	pb.RegisterProductServiceServer(grpcServer, productService)

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
	err = pb.RegisterProductServiceHandlerFromEndpoint(context.Background(), gwmux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway handler: %v", err)
	}

	// Start HTTP server for the gRPC-Gateway
	log.Println("gRPC-Gateway server running on port 8080")
	if err := http.ListenAndServe(":8080", gwmux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

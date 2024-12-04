package main

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"user-service/db"
	"user-service/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "user-service/pb"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func main() {
	// database connection
	dsn := "host=localhost user=postgres password=123 dbname=user_dev_db port=5432 sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// migrate models
	if err := gormDB.AutoMigrate(&pb.UserORM{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	dbProvider := db.NewAuthProvider(gormDB)

	// JWT Manager
	//jwtManager := jwt.NewJWTManager("mongodebol")
	//
	//// gRPC Server
	//interceptor := service.NewAuthInterceptor(jwtManager, service.AccessibleMethods())
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(interceptor.Unary()),
	//grpc.StreamInterceptor(interceptor.Stream()),
	)

	//register auth service
	authService := service.NewAuthServiceHandler(dbProvider)
	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("Serving gRPC on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// gRPC Gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err = pb.RegisterAuthServiceHandlerServer(ctx, mux, authService)
	if err != nil {
		log.Fatalf("Failed to register gRPC Gateway: %v", err)
	}

	log.Println("Serving HTTP on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

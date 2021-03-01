package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Bernigend/mb-cw3-phll-group-service/internal/app/endpoint"
	"github.com/Bernigend/mb-cw3-phll-group-service/internal/app/repository"
	"github.com/Bernigend/mb-cw3-phll-group-service/internal/app/service"

	api "github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api"
)

const (
	grpcServerPort = 8990
)

const (
	dbHost = "localhost"
	dbPort = "5434"
	dbUser = "user"
	dbPass = "password"
	dbName = "db"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := repository.NewRepository(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	log.Println("db ok")

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("migrations ok")

	srv := service.NewService(db)
	log.Println("service ok")

	endpoints := endpoint.NewEndpoint(srv)
	log.Println("endpoint ok")

	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServerPort))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen ok")

	grpcServer := grpc.NewServer()
	api.RegisterGroupServiceServer(grpcServer, endpoints)

	log.Printf("starting grpc server at 0.0.0.0:%v...\n", grpcServerPort)
	err = grpcServer.Serve(listenPort)
	if err != nil {
		log.Fatal(err)
	}
}

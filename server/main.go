package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yeaung276/ImageSearch/server/milvus"
	pb "github.com/yeaung276/ImageSearch/server/pb"
	"github.com/yeaung276/ImageSearch/server/service"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := milvus.New(ctx)
	if err != nil {
		log.Fatal("Error connecting to milvus server")
		return
	}
	defer db.Close()

	s := grpc.NewServer()
	srv, err := service.NewImageSearchServer(db)
	pb.RegisterImageSearchServiceServer(s, srv)
	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			log.Fatal("Error listening to port 50051")
			return
		}
		log.Print("Serving grpc server at port 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatal("Error serving grpc server")
			return
		}
	}()
	go func() {
		mux := runtime.NewServeMux()
		if err := pb.RegisterImageSearchServiceHandlerServer(ctx, mux, srv); err != nil {
			log.Fatal("Error registering Http handler")
			return
		}
		log.Print("Serving http server at port 50052")
		if err := http.ListenAndServe("0.0.0.0:50052", mux); err != nil {
			log.Fatal("Error serving http server")
			return
		}
	}()
	log.Print("Serving file server at port 9000")
	fs := http.FileServer(http.Dir("../jsmodel"))
	log.Fatal(http.ListenAndServe(":9000", fs))
}

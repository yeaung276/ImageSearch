package main

import (
	"context"
	"log"
	"net/http"

	"github.com/yeaung276/ImageSearch/server/service"
)

var opt service.Option

func init() {
	opt = service.Option{
		DBUrl:          "localhost:19530",
		CollectionName: "image_search",
		VectorSize:     2,
		GrpcAddress:    "0.0.0.0:50051",
		HTTPAddress:    "0.0.0.0:50052",
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := service.NewService(ctx, opt)
	defer svc.Shutdown(ctx)
	if err != nil {
		log.Fatal("Fail to start service", err)
		return
	}
	go func() {
		if err := svc.ServeGrpc(ctx); err != nil {
			log.Fatal("Fail to serve grpc server")
			return
		}
	}()
	go func() {
		if err := svc.ServeHttp(ctx); err != nil {
			log.Fatal("Fail to server http server")
		}
	}()
	log.Print("Serving file server at port 9000")
	fs := http.FileServer(http.Dir("../jsmodel"))
	log.Fatal(http.ListenAndServe(":9000", fs))
}

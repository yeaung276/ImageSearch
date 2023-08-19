package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/yeaung276/ImageSearch/server/service"
)

var opt service.Option

func init() {
	vs, _ := strconv.Atoi(os.Getenv("VECTOR_SIZE"))
	opt = service.Option{
		DBUrl:          os.Getenv("MILVUS_URL"),
		CollectionName: os.Getenv("COLLECTION_NAME"),
		VectorSize:     vs,
		GrpcAddress:    os.Getenv("GRPC_ADDRESS"),
		HTTPAddress:    os.Getenv("HTTP_ADDRESS"),
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

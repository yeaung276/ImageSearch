package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/yeaung276/ImageSearch/server/milvus"
	pb "github.com/yeaung276/ImageSearch/server/pb"
	"google.golang.org/grpc"
)

type Option struct {
	DBUrl          string
	CollectionName string
	VectorSize     int
	GrpcAddress    string
	HTTPAddress    string
}

type Service struct {
	GrpcAddress string
	HTTPAddress string
	srv         *ImageSearchServer
	db          *milvus.MilvusDB
}

func NewService(ctx context.Context, opt Option) (*Service, error) {
	db, err := milvus.New(ctx, milvus.Option{
		DBUrl:          opt.DBUrl,
		CollectionName: opt.CollectionName,
		VectorSize:     opt.VectorSize,
	})
	if err != nil {
		log.Fatal("Error connecting to milvus server")
		return nil, errors.New("Error connecting to milvus server")
	}

	srv, err := NewImageSearchServer(db, opt)
	return &Service{
		GrpcAddress: opt.GrpcAddress,
		HTTPAddress: opt.HTTPAddress,
		srv:         srv,
		db:          db,
	}, nil
}

func (s *Service) ServeGrpc(ctx context.Context) error {
	rpcs := grpc.NewServer()
	pb.RegisterImageSearchServiceServer(rpcs, s.srv)
	lis, err := net.Listen("tcp", s.GrpcAddress)
	if err != nil {
		log.Fatalf("Error listening to %s", s.GrpcAddress)
		return err
	}
	log.Printf("Serving grpc server at %s", s.GrpcAddress)
	if err := rpcs.Serve(lis); err != nil {
		log.Fatal("Error serving grpc server")
		return err
	}
	return nil
}

func (s *Service) ServeHttp(ctx context.Context) error {
	mux := runtime.NewServeMux()
	handler := cors.Default().Handler(mux)
	if err := pb.RegisterImageSearchServiceHandlerServer(ctx, mux, s.srv); err != nil {
		log.Fatal("Error registering Http handler")
		return err
	}
	log.Printf("Serving http server at %s", s.HTTPAddress)
	mux.HandlePath("POST", "/uploadFile", uploadFile)
	if err := http.ListenAndServe(s.HTTPAddress, handler); err != nil {
		log.Fatal("Error serving http server")
		return err
	}
	return nil
}

func (s *Service) Shutdown(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

// File upload
type fileResponse struct {
	FileUrl string
}

func uploadFile(w http.ResponseWriter, r *http.Request, params map[string]string) {
	//upload size
	err := r.ParseForm() // grab the multipart form
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("attachment")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file 'attachment': %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer f.Close()
	fileName := "/images/" + uuid.New().String()
	resFile, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer resFile.Close()
	if err == nil {
		io.Copy(resFile, f)
		defer resFile.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fileResponse{
			FileUrl: fileName,
		})
	}
}

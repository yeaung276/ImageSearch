package service

import (
	"context"
	"log"

	"github.com/yeaung276/ImageSearch/server/milvus"
	pb "github.com/yeaung276/ImageSearch/server/pb"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ImageSearchServer struct {
	pb.ImageSearchServiceServer
	db *milvus.MilvusDB
}

func NewImageSearchServer(db *milvus.MilvusDB) (*ImageSearchServer, error) {
	return &ImageSearchServer{
		db: db,
	}, nil
}

func (s *ImageSearchServer) Search(ctx context.Context, req *pb.ImageSearchRequest) (*pb.ImageSearchResponse, error) {
	if len(req.Embeddings) != 2 {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "Vector size mismatch")
	}
	r, err := s.db.SearchImages(ctx, req.Embeddings)
	if err != nil {
		log.Fatal("s.db.SearchImages: ", err)
		return nil, grpcstatus.Error(grpccodes.Internal, "Internal error")
	}

	result, err := milvus.NewResult(r)
	if err != nil {
		log.Fatal("milvus.NewResult", err)
		return nil, grpcstatus.Error(grpccodes.Internal, "Internal error")
	}

	return result, nil
}

func (s *ImageSearchServer) Add(ctx context.Context, req *pb.ImageAddRequest) (*emptypb.Empty, error) {
	if len(req.Embeddings) != 2 {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "Vector size mismatch")
	}

	if req.ImageUrl == "" {
		return nil, grpcstatus.Error(grpccodes.InvalidArgument, "ImageUrl is required")
	}

	if err := s.db.AddImageVector(ctx, req.ImageUrl, req.Embeddings); err != nil {
		log.Fatal("s.db.AddImageVector: ", err)
		return nil, grpcstatus.Error(grpccodes.Internal, "Error adding vector")
	}
	return nil, nil
}

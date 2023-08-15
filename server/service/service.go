package service

import (
	"context"
	"log"

	"github.com/yeaung276/ImageSearch/server/milvus"
	pb "github.com/yeaung276/ImageSearch/server/pb"
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

func (s *ImageSearchServer) Search(ctx context.Context, req *pb.ImageEmbedding) (*pb.ImageSearchResponse, error) {
	log.Print("Called from client")
	return &pb.ImageSearchResponse{}, nil
}

package milvus

import (
	"errors"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	pb "github.com/yeaung276/ImageSearch/server/pb"
)

func NewResult(sr []client.SearchResult) (*pb.ImageSearchResponse, error) {
	if len(sr) > 0 {
		r := &sr[0]
		var results []*pb.ImageResult
		field, ok := r.Fields[0].(*entity.ColumnVarChar)
		if !ok {
			return nil, errors.New("fail to typecast field.(*entity.ColumnVarChar).")
		}
		for i := 0; i < r.ResultCount; i++ {
			url, _ := field.ValueByIdx(i)
			results = append(results, &pb.ImageResult{
				Url:        url,
				Similarity: r.Scores[i],
			})
		}
		return &pb.ImageSearchResponse{
			Result:      results,
			ResultCount: int32(r.ResultCount),
		}, nil
	}
	return &pb.ImageSearchResponse{
		ResultCount: 0,
	}, nil
}

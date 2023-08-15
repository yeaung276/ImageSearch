package milvus

import (
	"context"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusDB struct {
	c client.Client
}

func New(ctx context.Context) (*MilvusDB, error) {
	mc, err := client.NewGrpcClient(
		ctx,
		"localhost:19530",
	)
	if err != nil {
		return nil, err
	}
	return &MilvusDB{
		c: mc,
	}, nil
}

func (m *MilvusDB) Close() error {
	return m.c.Close()
}

func (m *MilvusDB) CreateCollection(ctx context.Context) error {
	collExists, err := m.c.HasCollection(ctx, "test_collection")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if !collExists {
		schema := &entity.Schema{
			CollectionName: "test_collection",
			Description:    "this is the basic example collection",
			AutoID:         true,
			Fields: []*entity.Field{
				// currently primary key field is compulsory, and only int64 is allowd
				{
					Name:       "int64",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     true,
				},
				// also the vector field is needed
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{ // the vector dim may changed def method in release
						entity.TypeParamDim: "2",
					},
				},
			},
		}
		err = m.c.CreateCollection(ctx, schema, entity.DefaultShardNumber)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		idx, err := entity.NewIndexIvfFlat(entity.IP, 2)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		err = m.c.CreateIndex(ctx, "test_collection", "vector", idx, false)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
	}
	return nil
}

func (m *MilvusDB) DropCollection(ctx context.Context) error {
	collExists, err := m.c.HasCollection(ctx, "test_collection")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if collExists {
		return m.c.DropCollection(ctx, "test_collection")
	}
	return nil
}

func (m *MilvusDB) AddImageVector(ctx context.Context) error {
	vects := make([][]float32, 0, 1)
	vec := make([]float32, 2, 2)
	vects = append(vects, vec)
	vectorColumn := entity.NewColumnFloatVector("vector", 2, vects)
	_, err := m.c.Insert(ctx, "test_collection", "", vectorColumn)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (m *MilvusDB) ListImages(ctx context.Context) error {
	collExists, err := m.c.HasCollection(ctx, "test_collection")
	log.Print(collExists)
	err = m.c.LoadCollection(
		context.Background(), // ctx
		"test_collection",    // CollectionName
		false,                // async
	)
	if err != nil {
		log.Fatal("failed to load collection:", err.Error())
		return err
	}
	opt := client.SearchQueryOptionFunc(func(option *client.SearchQueryOption) {
		option.Limit = 3
		option.Offset = 0
		option.ConsistencyLevel = entity.ClStrong
		option.IgnoreGrowing = false
	})
	queryResult, err := m.c.Query(
		context.Background(), // ctx
		"test_collection",    // CollectionName
		[]string{},           // PartitionName
		"1==1",               // expr
		[]string{"int64"},    // OutputFields
		opt,                  // queryOptions
	)
	if err != nil {
		log.Fatal("fail to query collection:", err.Error())
		return err
	}
	log.Print(queryResult)
	return nil
}

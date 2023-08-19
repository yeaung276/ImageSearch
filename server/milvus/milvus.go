package milvus

import (
	"context"
	"log"
	"strconv"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Option struct {
	DBUrl          string
	CollectionName string
	VectorSize     int
}

type MilvusDB struct {
	c              client.Client
	collectionName string
	vectorSize     int
}

func New(ctx context.Context, opt Option) (*MilvusDB, error) {
	log.Printf("Connecting to malvus server at %s", opt.DBUrl)
	mc, err := client.NewGrpcClient(
		ctx,
		opt.DBUrl,
	)
	if err != nil {
		log.Fatal("Cannot start milvus grpc client")
		return nil, err
	}

	db := &MilvusDB{
		c:              mc,
		collectionName: opt.CollectionName,
		vectorSize:     opt.VectorSize,
	}

	if err := db.Setup(ctx); err != nil {
		log.Fatal("Milvus setup failed", err)
		return nil, err
	}

	return db, nil
}

func (m *MilvusDB) Close() error {
	return m.c.Close()
}

func (m *MilvusDB) Setup(ctx context.Context) error {
	log.Print("Setting up vector database...")
	collExists, err := m.c.HasCollection(ctx, m.collectionName)
	if err != nil {
		log.Fatal("Check collection exist failed", err.Error())
		return err
	}
	log.Print("Collection exist, skipping...")
	if !collExists {
		schema := &entity.Schema{
			CollectionName: m.collectionName,
			Description:    "A collection to store image embeddings",
			AutoID:         true,
			Fields: []*entity.Field{
				// currently primary key field is compulsory, and only int64 is allowd
				{
					Name:       "ID",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     true,
				},
				// currently primary key field is compulsory, and only int64 is allowd
				{
					Name:     "url",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						"max_length": "100",
					},
				},
				// also the vector field is needed
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{ // the vector dim may changed def method in release
						entity.TypeParamDim: strconv.Itoa(m.vectorSize),
					},
				},
			},
		}
		log.Print("Creating collection...")
		err = m.c.CreateCollection(ctx, schema, entity.DefaultShardNumber)
		if err != nil {
			log.Fatal("Creating collection failed. ", err.Error())
			return err
		}
		idx, err := entity.NewIndexIvfFlat(entity.IP, 2)
		if err != nil {
			log.Fatal("Creating index failed. ", err.Error())
			return err
		}
		log.Print("Creating index...")
		err = m.c.CreateIndex(ctx, m.collectionName, "vector", idx, false)
		if err != nil {
			log.Fatal("Creating index failed. ", err.Error())
			return err
		}
	}
	log.Print("db setup complete.")
	return nil
}

func (m *MilvusDB) DropCollection(ctx context.Context) error {
	collExists, err := m.c.HasCollection(ctx, m.collectionName)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if collExists {
		return m.c.DropCollection(ctx, m.collectionName)
	}
	return nil
}

func (m *MilvusDB) AddImageVector(ctx context.Context, url string, vector []float32) error {
	vects := make([][]float32, 0, 1)
	vects = append(vects, vector)
	vectorColumn := entity.NewColumnFloatVector("vector", 2, vects)
	urlColumn := entity.NewColumnVarChar("url", []string{url})
	_, err := m.c.Insert(ctx, m.collectionName, "", vectorColumn, urlColumn)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (m *MilvusDB) SearchImages(ctx context.Context, vector []float32) ([]client.SearchResult, error) {
	_, err := m.c.HasCollection(ctx, m.collectionName)
	if err != nil {
		log.Fatal("m.c.HasCollection()", err.Error())
		return nil, err
	}
	err = m.c.LoadCollection(
		ctx,              // ctx
		m.collectionName, // CollectionName
		false,            // async
	)
	if err != nil {
		log.Fatal("m.c.LoadCollection()", err.Error())
		return nil, err
	}
	opt := client.SearchQueryOptionFunc(func(option *client.SearchQueryOption) {
		option.Limit = 100
		option.Offset = 0
		option.ConsistencyLevel = entity.ClStrong
		option.IgnoreGrowing = false
	})
	sp, _ := entity.NewIndexIvfFlatSearchParam( // NewIndex*SearchParam func
		10, // searchParam
	)
	searchResult, err := m.c.Search(
		ctx,              // ctx
		m.collectionName, // CollectionName
		[]string{},       // PartitionName
		"",               // expr
		[]string{"url"},  // OutputFields
		[]entity.Vector{entity.FloatVector(vector)}, // Vectors
		"vector",  // Vector field
		entity.IP, // Metric type
		10,        // Top K
		sp,        // Search param
		opt,       // queryOptions
	)
	if err != nil {
		log.Fatal("m.c.Query()", err.Error())
		return nil, err
	}
	return searchResult, nil
}

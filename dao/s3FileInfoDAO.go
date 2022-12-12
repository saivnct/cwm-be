package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"sol.go/cwm/model"
	"sync"
	"time"
)

type S3FileInfoDAO struct {
	DAO
}

var singletonS3FileInfoDAO *S3FileInfoDAO
var onceS3FileInfoDAO sync.Once

func GetS3FileInfoDAO() *S3FileInfoDAO {
	onceS3FileInfoDAO.Do(func() {
		fmt.Println("Init S3FileInfoDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		s3FileInfoDAO := S3FileInfoDAO{}
		s3FileInfoDAO.Init(mongoCtx, &db.MongoDb)

		singletonS3FileInfoDAO = &s3FileInfoDAO
	})
	return singletonS3FileInfoDAO
}

func (s3FileInfoDAO *S3FileInfoDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "s3FileInfo"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	s3FileInfoDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (s3FileInfoDAO *S3FileInfoDAO) Save(ctx context.Context, s3FileInfo *model.S3FileInfo) (*model.S3FileInfo, error) {
	result := &model.S3FileInfo{}
	err := s3FileInfoDAO.InsertOrUpdate(ctx, s3FileInfo, result)
	if err != nil {
		return nil, fmt.Errorf("(S3FileInfoDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (s3FileInfoDAO *S3FileInfoDAO) FindByFileId(ctx context.Context, fileId string) (*model.S3FileInfo, error) {
	result := &model.S3FileInfo{}

	err := s3FileInfoDAO.FindByPKey(ctx, fileId, result)
	if err != nil {
		return nil, fmt.Errorf("(S3FileInfoDAO - FindByFileId): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (s3FileInfoDAO *S3FileInfoDAO) DeleteByFileId(ctx context.Context, fileId string) (*model.S3FileInfo, error) {
	result := &model.S3FileInfo{}
	err := s3FileInfoDAO.DeleteByPKey(ctx, fileId, result)
	if err != nil {
		return nil, fmt.Errorf("(S3FileInfoDAO - DeleteByFileId): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

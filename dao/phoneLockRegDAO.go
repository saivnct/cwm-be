package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"sol.go/cwm/model"
	"sync"
	"time"
)

type PhoneLockRegDAO struct {
	DAO
}

var singletonPhoneLockRegDAO *PhoneLockRegDAO
var oncePhoneLockRegDAO sync.Once

func GetPhoneLockRegDAO() *PhoneLockRegDAO {
	oncePhoneLockRegDAO.Do(func() {
		fmt.Println("Init PhoneLockRegDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		phoneLockRegDAO := PhoneLockRegDAO{}
		phoneLockRegDAO.Init(mongoCtx, &db.MongoDb)

		singletonPhoneLockRegDAO = &phoneLockRegDAO
	})
	return singletonPhoneLockRegDAO
}

func (phoneLockRegDAO *PhoneLockRegDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "phoneLockRegs"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	phoneLockRegDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (phoneLockRegDAO *PhoneLockRegDAO) Save(ctx context.Context, phoneLockReg *model.PhoneLockReg) (*model.PhoneLockReg, error) {
	result := &model.PhoneLockReg{}
	err := phoneLockRegDAO.InsertOrUpdate(ctx, phoneLockReg, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneLockRegDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (phoneLockRegDAO *PhoneLockRegDAO) FindByPhoneFull(ctx context.Context, phoneFull string) (*model.PhoneLockReg, error) {
	result := &model.PhoneLockReg{}

	err := phoneLockRegDAO.FindByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneLockRegDAO - FindByPhoneFull): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (phoneLockRegDAO *PhoneLockRegDAO) DeleteByPhoneFull(ctx context.Context, phoneFull string) (*model.PhoneLockReg, error) {
	result := &model.PhoneLockReg{}
	err := phoneLockRegDAO.DeleteByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneLockRegDAO - DeleteByPhoneFull): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

func (phoneLockRegDAO *PhoneLockRegDAO) UpdateByPhoneFull(ctx context.Context, phoneFull string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.PhoneLockReg, error) {
	result := &model.PhoneLockReg{}

	err := phoneLockRegDAO.UpdateByPKey(ctx, phoneFull, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneLockRegDAO - UpdateByPhoneFull): failed executing UpdateOneByField -> %w", err)
	}
	return result, nil
}

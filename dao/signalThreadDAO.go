package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"log"
	"sol.go/cwm/model"
	"sync"
	"time"
)

type SignalThreadDAO struct {
	DAO
}

var singletonSignalThreadDAO *SignalThreadDAO
var onceSignalThreadDAO sync.Once

func GetSignalThreadDAO() *SignalThreadDAO {
	onceSignalThreadDAO.Do(func() {
		fmt.Println("Init UserContactDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		signalThreadDAO := SignalThreadDAO{}
		signalThreadDAO.Init(mongoCtx, &db.MongoDb)

		singletonSignalThreadDAO = &signalThreadDAO
	})
	return singletonSignalThreadDAO
}

func (signalThreadDAO *SignalThreadDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "signalThreads"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	signalThreadDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (signalThreadDAO *SignalThreadDAO) Save(ctx context.Context, signalThread *model.SignalThread) (*model.SignalThread, error) {
	result := &model.SignalThread{}
	err := signalThreadDAO.InsertOrUpdate(ctx, signalThread, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (signalThreadDAO *SignalThreadDAO) FindByThreadId(ctx context.Context, threadId string) (*model.SignalThread, error) {
	result := &model.SignalThread{}

	err := signalThreadDAO.FindByPKey(ctx, threadId, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - FindByThreadId): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (signalThreadDAO *SignalThreadDAO) DeleteByThreadId(ctx context.Context, threadId string) (*model.SignalThread, error) {
	result := &model.SignalThread{}
	err := signalThreadDAO.DeleteByPKey(ctx, threadId, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - DeleteByThreadId): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

func (signalThreadDAO *SignalThreadDAO) UpdateByThreadId(ctx context.Context, threadId string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.SignalThread, error) {
	result := &model.SignalThread{}

	err := signalThreadDAO.UpdateByPKey(ctx, threadId, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - UpdateByPhoneFull): failed executing UpdateOneByField -> %w", err)
	}
	return result, nil
}

func (signalThreadDAO *SignalThreadDAO) UpdateParticipants(ctx context.Context, threadId string, participantAdds []string, participantRemoves []string) (*model.SignalThread, error) {
	mutexName := fmt.Sprintf("%v_%v", signalThreadDAO.CollectionName, threadId)

	redLock := signalThreadDAO.CreateRedlock(ctx, mutexName, signalThreadDAO.CacheLockTTL)

	// Obtain a lock for our given mutex. After this is successful, no one else
	//	// can obtain the same lock (the same mutex name) until we unlock it.
	log.Println("try acquire Lock", mutexName)
	err := redLock.Lock()
	// Release the lock so other processes or threads can obtain a lock.
	defer func() {
		//redLock.Unlock()
		ok, err := redLock.Unlock()
		log.Println("release Lock", mutexName, ok, err)
	}()

	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - UpdateParticipants): failed executing redLock.Lock -> %w", err)
	}

	signalThread, err := signalThreadDAO.FindByThreadId(ctx, threadId)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - UpdateParticipants): failed executing UpdateParticipants -> %w", err)
	}

	updateFields := primitive.M{}

	participants := signalThread.Participants
	if len(participantAdds) > 0 {
		participants = append(participants, participantAdds...)

		allParticipants := signalThread.AllParticipants
		allParticipants = append(allParticipants, participantAdds...)
		updateFields["allParticipants"] = allParticipants
	}

	if len(participantRemoves) > 0 {
		for _, participant := range participantRemoves {
			idx := slices.IndexFunc(participants, func(p string) bool { return p == participant })
			if idx >= 0 {
				participants = slices.Delete(participants, idx, idx+1)
			}
		}
	}

	updateFields["participants"] = participants
	updateFields["lastModified"] = time.Now().UnixMilli()
	update := primitive.M{"$set": updateFields}

	return signalThreadDAO.UpdateByThreadId(ctx, threadId, update, []interface{}{}, false)
}

func (signalThreadDAO *SignalThreadDAO) UpdateAdmin(ctx context.Context, threadId string, adminAdds []string, adminRemoves []string) (*model.SignalThread, error) {
	mutexName := fmt.Sprintf("%v_%v", signalThreadDAO.CollectionName, threadId)

	redLock := signalThreadDAO.CreateRedlock(ctx, mutexName, signalThreadDAO.CacheLockTTL)

	// Obtain a lock for our given mutex. After this is successful, no one else
	//	// can obtain the same lock (the same mutex name) until we unlock it.
	log.Println("try acquire Lock", mutexName)
	err := redLock.Lock()
	// Release the lock so other processes or threads can obtain a lock.
	defer func() {
		//redLock.Unlock()
		ok, err := redLock.Unlock()
		log.Println("release Lock", mutexName, ok, err)
	}()

	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - UpdateParticipants): failed executing redLock.Lock -> %w", err)
	}

	signalThread, err := signalThreadDAO.FindByThreadId(ctx, threadId)
	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - UpdateParticipants): failed executing UpdateParticipants -> %w", err)
	}

	admins := signalThread.Admins
	if len(adminAdds) > 0 {
		admins = append(admins, adminAdds...)
	}

	if len(adminRemoves) > 0 {
		for _, admin := range adminRemoves {
			idx := slices.IndexFunc(admins, func(p string) bool { return p == admin })
			if idx >= 0 {
				admins = slices.Delete(admins, idx, idx+1)
			}
		}
	}

	updateFields := primitive.M{}
	updateFields["admins"] = admins
	updateFields["lastModified"] = time.Now().UnixMilli()
	update := primitive.M{"$set": updateFields}

	return signalThreadDAO.UpdateByThreadId(ctx, threadId, update, []interface{}{}, false)
}

func (signalThreadDAO *SignalThreadDAO) FindAllThreadOfUser(ctx context.Context, phoneFull string) ([]*model.SignalThread, error) {
	results := []*model.SignalThread{}

	filter := primitive.M{
		"participants": phoneFull,
	}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})

	err := signalThreadDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalThreadDAO - FindAllThreadOfUser): failed executing FindByField -> %w", err)
	}

	return results, nil
}

//func (signalThreadDAO *SignalThreadDAO) UpdateIncreaseSeqByThreadId(ctx context.Context, threadId string) (*model.SignalThread, error) {
//	result := &model.SignalThread{}
//
//	updateFields := primitive.M{}
//	updateFields["currentSeq"] = 1
//	update := primitive.M{"$inc": updateFields}
//
//	err := signalThreadDAO.UpdateByPKey(ctx, threadId, update, []interface{}{}, false, result)
//	if err != nil {
//		return nil, fmt.Errorf("(SignalThreadDAO - IncreaseSeqByThreadId): failed executing UpdateOneByField -> %w", err)
//	}
//	return result, nil
//}

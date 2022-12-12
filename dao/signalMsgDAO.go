package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sol.go/cwm/model"
	"sync"
	"time"
)

type SignalMsgDAO struct {
	DAO
}

var singletonSignalMsgDAO *SignalMsgDAO
var onceSignalMsgDAO sync.Once

func GetSignalMsgDAO() *SignalMsgDAO {
	onceSignalMsgDAO.Do(func() {
		fmt.Println("Init SignalMsgDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		signalMsgDAO := SignalMsgDAO{}
		signalMsgDAO.Init(mongoCtx, &db.MongoDb)

		singletonSignalMsgDAO = &signalMsgDAO
	})
	return singletonSignalMsgDAO
}

func (signalMsgDAO *SignalMsgDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "signalMsgs"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	signalMsgDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (signalMsgDAO *SignalMsgDAO) Save(ctx context.Context, signalMsg *model.SignalMsg) (*model.SignalMsg, error) {
	result := &model.SignalMsg{}
	err := signalMsgDAO.InsertOrUpdate(ctx, signalMsg, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (signalMsgDAO *SignalMsgDAO) FindByMsgId(ctx context.Context, msgId string) (*model.SignalMsg, error) {
	result := &model.SignalMsg{}

	err := signalMsgDAO.FindByPKey(ctx, msgId, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindByMsgId): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (signalMsgDAO *SignalMsgDAO) DeleteByMsgId(ctx context.Context, msgId string) (*model.SignalMsg, error) {
	result := &model.SignalMsg{}
	err := signalMsgDAO.DeleteByPKey(ctx, msgId, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - DeleteByMsgId): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

func (signalMsgDAO *SignalMsgDAO) UpdateByMsgId(ctx context.Context, msgId string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.SignalMsg, error) {
	result := &model.SignalMsg{}

	err := signalMsgDAO.UpdateByPKey(ctx, msgId, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - UpdateByMsgId): failed executing UpdateOneByField -> %w", err)
	}
	return result, nil
}

func (signalMsgDAO *SignalMsgDAO) FindAllByThreadId(ctx context.Context, threadId string, fromSeq *int64) ([]*model.SignalMsg, error) {
	results := []*model.SignalMsg{}

	filter := primitive.M{}
	filter["threadId"] = threadId
	if fromSeq != nil {
		filter["seq"] = primitive.M{"$gte": fromSeq}
	}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"seq": 1})

	err := signalMsgDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindAllByThreadId): failed executing FindByField -> %w", err)
	}

	return results, nil
}

func (signalMsgDAO *SignalMsgDAO) CountByThreadId(ctx context.Context, threadId string, fromSeq *int64) (int64, error) {

	filter := primitive.M{}
	filter["threadId"] = threadId
	if fromSeq != nil {
		filter["seq"] = primitive.M{"$gte": fromSeq}
	}

	return signalMsgDAO.CountDocuments(ctx, filter)
}

func (signalMsgDAO *SignalMsgDAO) UnreceivedFilter(phoneFull string, session string, userGroupThreads []string, fromDate int64, toDate int64) interface{} {
	/*
		{
			$and: [
				$or: [
					"from": phoneFull,
					"to": phoneFull,
					"threadId": {
						$in: threadIds
					}
				],
				"deleteForUsers": {
					$nin: []string{phoneFull}
				},
				"receivedSessions": {
					"$nin": []string{session},
				},
				"createdAt" : {
					"$gte": fromDate,
				},
				"createdAt" : {
					"$lte": toDate,
				},
			]
		}
	*/

	filter := primitive.M{}

	orArr1 := []interface{}{
		primitive.M{"from": phoneFull},
		primitive.M{"to": phoneFull},
	}
	if len(userGroupThreads) > 0 {
		orArr1 = append(orArr1, primitive.M{
			"threadId": primitive.M{
				"$in": userGroupThreads,
			},
		})
	}

	andElement1 := primitive.M{"$or": orArr1}

	andElement2 := primitive.M{
		"deleteForUsers": primitive.M{
			"$nin": []string{phoneFull},
		},
	}
	andElement3 := primitive.M{
		"receivedSessions": primitive.M{
			"$nin": []string{session},
		},
	}

	filterAnds := []interface{}{andElement1, andElement2, andElement3}
	if fromDate > 0 {
		filterAnds = append(filterAnds, primitive.M{
			"createdAt": primitive.M{
				"$gt": fromDate,
			},
		})
	}
	if toDate > 0 {
		filterAnds = append(filterAnds, primitive.M{
			"createdAt": primitive.M{
				"$lt": toDate,
			},
		})
	}

	filter["$and"] = filterAnds

	return filter
}

func (signalMsgDAO *SignalMsgDAO) FindAllUnreceived(ctx context.Context, phoneFull string, session string, userGroupThreads []string, fromDate int64, toDate int64, page, pageSize int64) ([]*model.SignalMsg, error) {
	results := []*model.SignalMsg{}

	filter := signalMsgDAO.UnreceivedFilter(phoneFull, session, userGroupThreads, fromDate, toDate)

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})
	if pageSize > 0 {
		findOptions.SetSkip(page * pageSize)
		findOptions.SetLimit(pageSize)
	}

	err := signalMsgDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindAllUnreceived): failed executing FindByField -> %w", err)
	}

	return results, nil
}

func (signalMsgDAO *SignalMsgDAO) CountAllUnreceived(ctx context.Context, phoneFull string, session string, userGroupThreads []string, fromDate int64, toDate int64) (int64, error) {
	filter := signalMsgDAO.UnreceivedFilter(phoneFull, session, userGroupThreads, fromDate, toDate)

	return signalMsgDAO.CountDocuments(ctx, filter)
}

func (signalMsgDAO *SignalMsgDAO) UnreceivedOfThreadFilter(phoneFull string, session string, threadId string, fromDate int64, toDate int64) interface{} {
	/*
		{
			$and: [
				"threadId": threadId,
				"deleteForUsers": {
					$nin: []string{phoneFull}
				},
				"receivedSessions": {
					"$nin": []string{session},
				},
				"createdAt" : {
					"$gte": fromDate,
				},
				"createdAt" : {
					"$lte": toDate,
				},
			]
		}
	*/

	//NOTE : Must check user is thread's member

	filter := primitive.M{}

	andElement1 := primitive.M{"threadId": threadId}

	andElement2 := primitive.M{
		"deleteForUsers": primitive.M{
			"$nin": []string{phoneFull},
		},
	}
	andElement3 := primitive.M{
		"receivedSessions": primitive.M{
			"$nin": []string{session},
		},
	}

	filterAnds := []interface{}{andElement1, andElement2, andElement3}
	if fromDate > 0 {
		filterAnds = append(filterAnds, primitive.M{
			"createdAt": primitive.M{
				"$gte": fromDate,
			},
		})
	}
	if toDate > 0 {
		filterAnds = append(filterAnds, primitive.M{
			"createdAt": primitive.M{
				"$lte": toDate,
			},
		})
	}

	filter["$and"] = filterAnds

	return filter
}

func (signalMsgDAO *SignalMsgDAO) FindAllUnreceivedOfThread(ctx context.Context, phoneFull string, session string,
	threadId string, fromDate int64, toDate int64,
	page, pageSize int64, sortOrder int) ([]*model.SignalMsg, error) {
	results := []*model.SignalMsg{}

	filter := signalMsgDAO.UnreceivedOfThreadFilter(phoneFull, session, threadId, fromDate, toDate)

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": sortOrder})
	if pageSize > 0 {
		findOptions.SetSkip(page * pageSize)
		findOptions.SetLimit(pageSize)
	}

	err := signalMsgDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindAllUnreceivedOfThread): failed executing FindByField -> %w", err)
	}

	return results, nil
}

func (signalMsgDAO *SignalMsgDAO) CountAllUnreceivedOfThread(ctx context.Context, phoneFull string, session string, threadId string, fromDate int64, toDate int64) (int64, error) {
	filter := signalMsgDAO.UnreceivedOfThreadFilter(phoneFull, session, threadId, fromDate, toDate)

	return signalMsgDAO.CountDocuments(ctx, filter)
}

func (signalMsgDAO *SignalMsgDAO) OldMsgOfThreadFilter(phoneFull string, threadId string, toDate int64) interface{} {
	/*
		{
			$and: [
				"threadId": threadId,
				"deleteForUsers": {
					$nin: []string{phoneFull}
				},
				"createdAt" : {
					"$lte": toDate,
				},
			]
		}
	*/

	//NOTE : Must check user is thread's member

	filter := primitive.M{}

	andElement1 := primitive.M{"threadId": threadId}

	andElement2 := primitive.M{
		"deleteForUsers": primitive.M{
			"$nin": []string{phoneFull},
		},
	}

	filterAnds := []interface{}{andElement1, andElement2}
	if toDate > 0 {
		filterAnds = append(filterAnds, primitive.M{
			"createdAt": primitive.M{
				"$lt": toDate,
			},
		})
	}

	filter["$and"] = filterAnds

	return filter
}

func (signalMsgDAO *SignalMsgDAO) FindOldMsgOfThread(ctx context.Context,
	phoneFull string, threadId string, toDate int64, limit int64) ([]*model.SignalMsg, error) {

	results := []*model.SignalMsg{}

	filter := signalMsgDAO.OldMsgOfThreadFilter(phoneFull, threadId, toDate)

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": -1}) //Sort descending by createdAt
	findOptions.SetLimit(limit)

	err := signalMsgDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindOldMsgOfThread): failed executing FindByField -> %w", err)
	}

	return results, nil
}

func (signalMsgDAO *SignalMsgDAO) AppendReceivedSessionByMsgIds(ctx context.Context, msgIds []string, sessionId string) (int64, error) {
	//https://www.mongodb.com/docs/manual/reference/operator/update/addToSet/
	filter := primitive.M{
		"pkey": primitive.M{
			"$in": msgIds,
		},
	}

	updateFields := primitive.M{}
	updateFields["receivedSessions"] = sessionId
	update := primitive.M{"$addToSet": updateFields}

	return signalMsgDAO.UpdateMany(ctx, filter, update)
}

func (signalMsgDAO *SignalMsgDAO) FindAllByMsgIds(ctx context.Context, msgIds []string) ([]*model.SignalMsg, error) {
	results := []*model.SignalMsg{}
	filter := primitive.M{
		"pkey": primitive.M{
			"$in": msgIds,
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})

	err := signalMsgDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(SignalMsgDAO - FindAllUnreceived): failed executing FindByField -> %w", err)
	}

	return results, nil
}

func (signalMsgDAO *SignalMsgDAO) AppendSeenByUserByMsgIds(ctx context.Context, msgIds []string, phoneFull string) (int64, error) {
	//https://www.mongodb.com/docs/manual/reference/operator/update/addToSet/
	filter := primitive.M{
		"pkey": primitive.M{
			"$in": msgIds,
		},
	}

	updateFields := primitive.M{}
	updateFields["seenByUsers"] = phoneFull
	update := primitive.M{"$addToSet": updateFields}

	return signalMsgDAO.UpdateMany(ctx, filter, update)
}

func (signalMsgDAO *SignalMsgDAO) AppendDeleteUsersByMsgIds(ctx context.Context, msgIds []string, phoneFulls []string) (int64, error) {
	//https://www.mongodb.com/docs/manual/reference/operator/update/addToSet/
	filter := primitive.M{
		"pkey": primitive.M{
			"$in": msgIds,
		},
	}

	updateFields := primitive.M{}
	updateFields["deleteForUsers"] = primitive.M{
		"$each": phoneFulls,
	}
	update := primitive.M{"$addToSet": updateFields}

	return signalMsgDAO.UpdateMany(ctx, filter, update)
}

func (signalMsgDAO *SignalMsgDAO) AppendDeleteUsersByThreadId(ctx context.Context, threadId string, phoneFulls []string) (int64, error) {
	//https://www.mongodb.com/docs/manual/reference/operator/update/addToSet/
	filter := primitive.M{
		"threadId": threadId,
	}

	updateFields := primitive.M{}
	updateFields["deleteForUsers"] = primitive.M{
		"$each": phoneFulls,
	}
	update := primitive.M{"$addToSet": updateFields}

	return signalMsgDAO.UpdateMany(ctx, filter, update)
}

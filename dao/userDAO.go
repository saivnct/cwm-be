package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"log"
	"sol.go/cwm/model"
	"sol.go/cwm/utils"
	"sync"
	"time"
)

const (
	User_FIELD_USERNAME = "username"
)

type UserDAO struct {
	DAO
}

var singletonUserDAO *UserDAO
var onceUserDAO sync.Once

func GetUserDAO() *UserDAO {
	onceUserDAO.Do(func() {
		fmt.Println("Init UserDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		userDAO := UserDAO{}
		userDAO.Init(mongoCtx, &db.MongoDb)

		singletonUserDAO = &userDAO
	})
	return singletonUserDAO
}

func (userDAO *UserDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "users"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	userDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (userDAO *UserDAO) Save(ctx context.Context, user *model.User) (*model.User, error) {
	result := &model.User{}
	err := userDAO.InsertOrUpdate(ctx, user, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - Save): failed executing Save -> %w", err)
	}

	err = userDAO.SaveToCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - Save): failed executing SaveToCache -> %w", err)
	}

	return result, nil
}

func (userDAO *UserDAO) FindByUserName(ctx context.Context, username string) (*model.User, error) {
	result := &model.User{}

	err := userDAO.GetFromCache(ctx, utils.CacheKey(userDAO.CollectionName, User_FIELD_USERNAME, username), result)
	if err == nil {
		log.Println("FindByUserName found in cache OK")
		return result, nil
	}

	log.Println("FindByUserName not found in cache -> find in db")
	err = userDAO.FindByField(ctx, User_FIELD_USERNAME, username, result)

	if err != nil {
		return nil, fmt.Errorf("(UserDAO - FindByUserName): failed executing FindByField -> %w", err)
	}

	err = userDAO.SaveToCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - FindByUserName): failed executing SaveToCache -> %w", err)
	}

	return result, nil
}

func (userDAO *UserDAO) SearchByUserName(ctx context.Context, username string) ([]*model.User, error) {
	results := []*model.User{}

	//{username: {$regex: ".*GdK.*", $options: "i"}}

	pattern := fmt.Sprintf(".*%v.*", username)
	search := bson.M{
		"$regex":   pattern,
		"$options": "i",
	}
	filter := bson.M{User_FIELD_USERNAME: search}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})
	findOptions.SetLimit(20)

	err := userDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(UserDAO - SearchByUserName): failed executing SearchByField -> %w", err)
	}

	return results, nil
}

func (userDAO *UserDAO) SearchByPhoneFull(ctx context.Context, phoneFull string) ([]*model.User, error) {
	results := []*model.User{}

	//{pkey: {$regex: ".*123456.*", $options: "i"}}

	pattern := fmt.Sprintf(".*%v.*", phoneFull)
	search := bson.M{
		"$regex":   pattern,
		"$options": "i",
	}
	filter := bson.M{PKEY_NAME: search}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})
	findOptions.SetLimit(20)

	err := userDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(UserDAO - SearchByUserName): failed executing SearchByField -> %w", err)
	}

	return results, nil
}

func (userDAO *UserDAO) FindByListPhoneFull(ctx context.Context, phoneFulls []string) ([]*model.User, error) {
	results := []*model.User{}

	//{pkey: {$in: ["+84966229268","+84123456789","+84888888888"]}}

	filter := bson.M{PKEY_NAME: primitive.M{
		"$in": phoneFulls,
	}}

	findOptions := options.Find()
	findOptions.SetSort(primitive.M{"createdAt": 1})

	err := userDAO.FindAll(ctx, filter, findOptions, &results)

	if err != nil {
		return nil, fmt.Errorf("(UserDAO - FindByListPhoneFull): failed executing SearchByField -> %w", err)
	}

	return results, nil
}

func (userDAO *UserDAO) FindByPhoneFull(ctx context.Context, phoneFull string) (*model.User, error) {
	result := &model.User{}

	err := userDAO.GetFromCache(ctx, utils.CacheKey(userDAO.CollectionName, PKEY_NAME, phoneFull), result)
	if err == nil {
		//log.Println("FindByPhoneFull found in cache OK")
		return result, nil
	}

	//log.Println("FindByPhoneFull not found in cache -> find in db")
	err = userDAO.FindByPKey(ctx, phoneFull, result)

	if err != nil {
		return nil, fmt.Errorf("(UserDAO - FindByPhoneFull): failed executing FindByField -> %w", err)
	}

	err = userDAO.SaveToCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - FindByPhoneFull): failed executing SaveToCache -> %w", err)
	}

	return result, nil
}

func (userDAO *UserDAO) DeleteByPhoneFull(ctx context.Context, phoneFull string) (*model.User, error) {
	result := &model.User{}
	err := userDAO.DeleteByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - DeleteByPhoneFull): failed executing DeleteOneByField -> %w", err)
	}

	err = userDAO.DeleteFromCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - DeleteByPhoneFull): failed executing DeleteFromCache -> %w", err)
	}
	return result, nil
}

func (userDAO *UserDAO) DeleteByUserName(ctx context.Context, username string) (*model.User, error) {
	result := &model.User{}
	err := userDAO.DeleteOneByField(ctx, User_FIELD_USERNAME, username, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - DeleteByUserName): failed executing DeleteOneByField -> %w", err)
	}

	err = userDAO.DeleteFromCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - DeleteByUserName): failed executing DeleteFromCache -> %w", err)
	}
	return result, nil
}

func (userDAO *UserDAO) UpdateByPhoneFull(ctx context.Context, phoneFull string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.User, error) {
	result := &model.User{}

	err := userDAO.UpdateByPKey(ctx, phoneFull, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateByPhoneFull): failed executing UpdateOneByField -> %w", err)
	}

	err = userDAO.SaveToCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateByPhoneFull): failed executing SaveToCache -> %w", err)
	}
	return result, nil
}

func (userDAO *UserDAO) UpdateByUserName(ctx context.Context, username string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.User, error) {
	result := &model.User{}

	err := userDAO.UpdateOneByField(ctx, User_FIELD_USERNAME, username, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateByUserName): failed executing UpdateOneByField -> %w", err)
	}

	err = userDAO.SaveToCache(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateByUserName): failed executing SaveToCache -> %w", err)
	}
	return result, nil
}

func (userDAO *UserDAO) SaveToCache(ctx context.Context, user *model.User) error {
	err := userDAO.SetToCache(ctx, utils.CacheKey(userDAO.CollectionName, PKEY_NAME, user.PhoneFull), user, userDAO.CacheTTL)
	if err != nil {
		return fmt.Errorf("(UserDAO - SaveToCache): failed executing SetToCache FIELD_EMAIL -> %w", err)
	}

	if len(user.Username) > 0 {
		err = userDAO.SetToCache(ctx, utils.CacheKey(userDAO.CollectionName, User_FIELD_USERNAME, user.Username), user, userDAO.CacheTTL)
		if err != nil {
			return fmt.Errorf("(UserDAO - SaveToCache): failed executing SetToCache User_FIELD_USERNAME -> %w", err)
		}
	}

	return nil
}

func (userDAO *UserDAO) DeleteFromCache(ctx context.Context, user *model.User) error {
	keys := []string{
		utils.CacheKey(userDAO.CollectionName, PKEY_NAME, user.PhoneFull),
		utils.CacheKey(userDAO.CollectionName, User_FIELD_USERNAME, user.Username),
	}
	err := userDAO.DelFromCache(ctx, keys)
	if err != nil {
		return fmt.Errorf("(UserDAO - DeleteFromCache): failed executing DelFromCache -> %w", err)
	}
	return nil
}

func (userDAO *UserDAO) UpdateNonce(ctx context.Context, phoneFull string) (*model.User, error) {
	mutexName := fmt.Sprintf("%v_%v", userDAO.CollectionName, phoneFull)

	redLock := userDAO.CreateRedlock(ctx, mutexName, userDAO.CacheLockTTL)

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
		return nil, fmt.Errorf("(UserDAO - UpdateNonce): failed executing redLock.Lock -> %w", err)
	}

	user, err := userDAO.FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateNonce): failed executing FindByPhoneFull -> %w", err)
	}

	if !user.IsExpireNonce() {
		return user, nil
	}

	shouldRenewNonce := user.ShouldRenewNonce()
	if !shouldRenewNonce {
		return user, nil
	}

	updateFields := primitive.M{}
	updateFields["nonce"] = user.Nonce
	updateFields["nonceResponse"] = user.NonceResponse
	updateFields["nonceAt"] = user.NonceAt
	update := primitive.M{"$set": updateFields}
	return userDAO.UpdateByPhoneFull(ctx, user.PhoneFull, update, []interface{}{}, false)
}

func (userDAO *UserDAO) UpdateGroupThread(ctx context.Context, phoneFull string, threadAdds []string, threadRemoves []string) (*model.User, error) {
	mutexName := fmt.Sprintf("%v_%v", userDAO.CollectionName, phoneFull)

	redLock := userDAO.CreateRedlock(ctx, mutexName, userDAO.CacheLockTTL)

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
		return nil, fmt.Errorf("(UserDAO - UpdateGroupThread): failed executing redLock.Lock -> %w", err)
	}

	user, err := userDAO.FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		return nil, fmt.Errorf("(UserDAO - UpdateGroupThread): failed executing FindByPhoneFull -> %w", err)
	}

	groupThreads := user.GroupThreads
	if len(threadAdds) > 0 {
		groupThreads = append(groupThreads, threadAdds...)
	}

	if len(threadRemoves) > 0 {
		for _, threadId := range threadRemoves {
			idx := slices.IndexFunc(groupThreads, func(t string) bool { return t == threadId })
			if idx >= 0 {
				groupThreads = slices.Delete(groupThreads, idx, idx+1)
			}
		}
	}

	updateFields := primitive.M{}
	updateFields["groupThreads"] = groupThreads
	update := primitive.M{"$set": updateFields}

	return userDAO.UpdateByPhoneFull(ctx, phoneFull, update, []interface{}{}, false)
}

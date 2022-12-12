package dao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"log"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/grpcCWMPb"
	"sync"
	"time"
)

type UserContactDAO struct {
	DAO
}

var singletonUserContactDAO *UserContactDAO
var onceUserContactDAO sync.Once

func GetUserContactDAO() *UserContactDAO {
	onceUserContactDAO.Do(func() {
		fmt.Println("Init UserContactDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		userContactDAO := UserContactDAO{}
		userContactDAO.Init(mongoCtx, &db.MongoDb)

		singletonUserContactDAO = &userContactDAO
	})
	return singletonUserContactDAO
}

func (userContactDAO *UserContactDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "userContacts"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	userContactDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (userContactDAO *UserContactDAO) Save(ctx context.Context, userContact *model.UserContact) (*model.UserContact, error) {
	result := &model.UserContact{}
	err := userContactDAO.InsertOrUpdate(ctx, userContact, result)
	if err != nil {
		return nil, fmt.Errorf("(UserContactDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (userContactDAO *UserContactDAO) FindByPhoneFull(ctx context.Context, phoneFull string) (*model.UserContact, error) {
	result := &model.UserContact{}

	err := userContactDAO.FindByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(UserContactDAO - FindByPhoneFull): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (userContactDAO *UserContactDAO) DeleteByPhoneFull(ctx context.Context, phoneFull string) (*model.UserContact, error) {
	result := &model.UserContact{}
	err := userContactDAO.DeleteByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(UserContactDAO - DeleteByPhoneFull): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

func (userContactDAO *UserContactDAO) UpdateByPhoneFull(ctx context.Context, phoneFull string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.UserContact, error) {
	result := &model.UserContact{}

	err := userContactDAO.UpdateByPKey(ctx, phoneFull, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(UserContactDAO - UpdateByPhoneFull): failed executing UpdateOneByField -> %w", err)
	}
	return result, nil
}

func (userContactDAO *UserContactDAO) UpdateContactList(ctx context.Context, user *model.User, sessionId string, contactInfos []*grpcCWMPb.ContactInfo) (*model.UserContact, error) {
	mutexName := fmt.Sprintf("%v_%v", userContactDAO.CollectionName, user.PhoneFull)
	redLock := userContactDAO.CreateRedlock(ctx, mutexName, userContactDAO.CacheLockTTL)

	// Obtain a lock for our given mutex. After this is successful, no one else
	//	// can obtain the same lock (the same mutex name) until we unlock it.
	//log.Println("try acquire Lock", mutexName)
	err := redLock.Lock()
	// Release the lock so other processes or threads can obtain a lock.
	defer func() {
		redLock.Unlock()
		//ok, err := redLock.Unlock()
		//log.Println("release Lock", mutexName, ok, err)
	}()

	if err != nil {
		return nil, fmt.Errorf("(UserContactDAO - UpdateContactList): failed executing redLock.Lock -> %w", err)
	}

	userContact, err := userContactDAO.FindByPhoneFull(ctx, user.PhoneFull)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //phoneLockReg == nil
			userContact = &model.UserContact{
				PhoneFull: user.PhoneFull,
				Contacts:  make(map[string]model.Contact),
			}
		} else {
			log.Println("UpdateContactList err - ", err)
			return nil, fmt.Errorf("(UserContactDAO - UpdateContactList): failed -> %w", err)
		}
	}

	shouldSave := false
	for _, contactInfo := range contactInfos {
		phoneFull := contactInfo.GetPhoneFull()
		if val, ok := userContact.Contacts[phoneFull]; ok {
			idx := slices.IndexFunc(val.Sessions, func(c string) bool { return c == sessionId })
			if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_ADD || contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
				if idx < 0 {
					val.Sessions = append(val.Sessions, sessionId)
					shouldSave = true
				}

				if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
					val.ContactName = contactInfo.GetName()
					val.PhoneFull = phoneFull
					shouldSave = true
				}

			} else if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_REMOVE {
				if idx >= 0 {
					val.Sessions = slices.Delete(val.Sessions, idx, idx+1)
					shouldSave = true
				}
			}
			userContact.Contacts[phoneFull] = val
		} else {
			if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_ADD || contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
				userContact.Contacts[phoneFull] = model.Contact{
					PhoneFull:   phoneFull,
					ContactName: contactInfo.GetName(),
					Sessions:    []string{sessionId},
				}
				shouldSave = true
			} else if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_REMOVE {
				userContact.Contacts[phoneFull] = model.Contact{
					PhoneFull:   phoneFull,
					ContactName: contactInfo.GetName(),
					Sessions:    []string{},
				}
				shouldSave = true
			}
		}
	}

	if shouldSave {
		//log.Println("Save update  userContact")
		return userContactDAO.Save(ctx, userContact)
	}

	//log.Println("Not Save update  userContact")
	return userContact, nil
}

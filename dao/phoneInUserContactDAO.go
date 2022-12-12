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

type PhoneInUserContactDAO struct {
	DAO
}

var singletonPhoneInUserContactDAO *PhoneInUserContactDAO
var oncePhoneInUserContactDAO sync.Once

func GetPhoneInUserContactDAO() *PhoneInUserContactDAO {
	oncePhoneInUserContactDAO.Do(func() {
		fmt.Println("Init PhoneInUserContactDAO...")

		db := GetDataBase()
		mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelMongo()

		phoneInUserContactDAO := PhoneInUserContactDAO{}
		phoneInUserContactDAO.Init(mongoCtx, &db.MongoDb)

		singletonPhoneInUserContactDAO = &phoneInUserContactDAO
	})
	return singletonPhoneInUserContactDAO
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) Init(ctx context.Context, db *mongo.Database) {
	COLLECTION_NAME := "phoneInUserContacts"
	CACHE_TTL := 10 * time.Minute
	CACHE_LOCK_TTL := 30 * time.Second
	phoneInUserContactDAO.InitDAO(ctx, db, COLLECTION_NAME, []string{}, CACHE_TTL, CACHE_LOCK_TTL)
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) Save(ctx context.Context, phoneInUserContact *model.PhoneInUserContact) (*model.PhoneInUserContact, error) {
	result := &model.PhoneInUserContact{}
	err := phoneInUserContactDAO.InsertOrUpdate(ctx, phoneInUserContact, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneInUserContactDAO - Save): failed executing Save -> %w", err)
	}

	return result, nil
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) FindByPhoneFull(ctx context.Context, phoneFull string) (*model.PhoneInUserContact, error) {
	result := &model.PhoneInUserContact{}

	err := phoneInUserContactDAO.FindByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneInUserContactDAO - FindByPhoneFull): failed executing FindByField -> %w", err)
	}

	return result, nil
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) DeleteByPhoneFull(ctx context.Context, phoneFull string) (*model.PhoneInUserContact, error) {
	result := &model.PhoneInUserContact{}
	err := phoneInUserContactDAO.DeleteByPKey(ctx, phoneFull, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneInUserContactDAO - DeleteByPhoneFull): failed executing DeleteOneByField -> %w", err)
	}
	return result, nil
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) UpdateByPhoneFull(ctx context.Context, phoneFull string, update interface{}, arrayFilter []interface{}, upsert bool) (*model.PhoneInUserContact, error) {
	result := &model.PhoneInUserContact{}

	err := phoneInUserContactDAO.UpdateByPKey(ctx, phoneFull, update, arrayFilter, upsert, result)
	if err != nil {
		return nil, fmt.Errorf("(PhoneInUserContactDAO - UpdateByPhoneFull): failed executing UpdateOneByField -> %w", err)
	}
	return result, nil
}

func (phoneInUserContactDAO *PhoneInUserContactDAO) UpdateUserList(ctx context.Context, contactInfo *grpcCWMPb.ContactInfo, user *model.User, sessionId string) (*model.PhoneInUserContact, error) {
	phoneFull := contactInfo.GetPhoneFull()
	mutexName := fmt.Sprintf("%v_%v", phoneInUserContactDAO.CollectionName, phoneFull)
	redLock := phoneInUserContactDAO.CreateRedlock(ctx, mutexName, phoneInUserContactDAO.CacheLockTTL)

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
		return nil, fmt.Errorf("(PhoneInUserContactDAO - UpdateUserList): failed executing redLock.Lock -> %w", err)
	}

	phoneInUserContact, err := phoneInUserContactDAO.FindByPhoneFull(ctx, phoneFull)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { //phoneLockReg == nil
			phoneInUserContact = &model.PhoneInUserContact{
				PhoneFull: phoneFull,
				Users:     make(map[string]model.ContactSession),
			}
		} else {
			log.Println("UpdateUserList err - ", err)
			return nil, fmt.Errorf("(PhoneInUserContactDAO - UpdateUserList): failed -> %w", err)
		}
	}

	shouldSave := false
	if val, ok := phoneInUserContact.Users[user.PhoneFull]; ok {
		idx := slices.IndexFunc(val.Sessions, func(c string) bool { return c == sessionId })
		if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_ADD || contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
			if idx < 0 {
				val.Sessions = append(val.Sessions, sessionId)
				phoneInUserContact.Users[user.PhoneFull] = val
				shouldSave = true
			}
		} else if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_REMOVE {
			if idx >= 0 {
				val.Sessions = slices.Delete(val.Sessions, idx, idx+1)
				phoneInUserContact.Users[user.PhoneFull] = val
				shouldSave = true
			}
		}

	} else {
		if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_ADD || contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
			phoneInUserContact.Users[user.PhoneFull] = model.ContactSession{
				PhoneFull: user.PhoneFull,
				Sessions:  []string{sessionId},
			}
			shouldSave = true
		} else if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_REMOVE {
			phoneInUserContact.Users[user.PhoneFull] = model.ContactSession{
				PhoneFull: user.PhoneFull,
				Sessions:  []string{},
			}
			shouldSave = true
		}
	}

	if shouldSave {
		//log.Println("Save update  phoneInUserContact")
		return phoneInUserContactDAO.Save(ctx, phoneInUserContact)
	}
	//log.Println("Not Save update  phoneInUserContact")
	return phoneInUserContact, nil
}

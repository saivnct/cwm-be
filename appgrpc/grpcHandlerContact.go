package appgrpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/proto/grpcCWMPb"
)

type ContactGRPCController struct {
}

func syncContactJob(contactInfoChannel chan *grpcCWMPb.ContactInfo, user *model.User, sessionId string, resultChannel chan string, stream grpcCWMPb.CWMService_SyncContactServer) {
	contactInfos := []*grpcCWMPb.ContactInfo{}
	for contactInfo := range contactInfoChannel {
		userId := ""
		userName := ""
		userAvatar := ""

		contactInfos = append(contactInfos, contactInfo)

		_, err := dao.GetPhoneInUserContactDAO().UpdateUserList(context.Background(), contactInfo, user, sessionId)

		if err != nil {
			log.Printf("Error while updating PhoneInUserContact: %v\n", err)
		} else {
			if contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_ADD || contactInfo.GetSyncType() == grpcCWMPb.CONTACT_SYNC_TYPE_UPDATE {
				user, _ := dao.GetUserDAO().FindByPhoneFull(context.Background(), contactInfo.GetPhoneFull())
				if user != nil {
					userId = user.ID.Hex()
					userName = user.Username
					userAvatar = user.Avatar
				}
			}

			resp := grpcCWMPb.SyncContactResponse{
				ContactInfo: &grpcCWMPb.ContactInfo{
					ContactId:  contactInfo.GetContactId(),
					Name:       contactInfo.GetName(),
					PhoneFull:  contactInfo.GetPhoneFull(),
					SyncType:   contactInfo.SyncType,
					UserId:     &userId,
					Username:   &userName,
					UserAvatar: &userAvatar,
				},
			}

			err = stream.Send(&resp)
			if err != nil {
				log.Printf("Error while sending response to client: %v\n", err)
			}
		}

	}

	_, err := dao.GetUserContactDAO().UpdateContactList(context.Background(), user, sessionId, contactInfos)
	if err != nil {
		log.Printf("Error while Update UserContact: %v\n", err)
	}

	resultChannel <- "Done"
	close(resultChannel)
}

func (sv *CWMGRPCService) SyncContact(stream grpcCWMPb.CWMService_SyncContactServer) error {

	grpcSession, ok := stream.Context().Value(GRPC_CTX_KEY_SESSION).(*GrpcSession)
	if !ok {
		log.Println("SyncContact - can not cast GrpcSession")
		return nil
	}

	resultChannel := make(chan string)
	contactInfoChannel := make(chan *grpcCWMPb.ContactInfo, 1000) //make a buffered channel size 1000
	go syncContactJob(contactInfoChannel, grpcSession.User, grpcSession.SessionId, resultChannel, stream)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			fmt.Println("Client stop sending stream")
			close(contactInfoChannel)
			break
		}

		if err != nil {
			log.Printf("Error while reading client stream: %v\n", err)
			close(contactInfoChannel)
			break
		}

		contactInfo := req.GetContactInfo()
		//log.Println("SyncContact", contactInfo)
		if len(contactInfo.GetName()) == 0 {
			continue
		}

		contactInfoChannel <- contactInfo
	}

	result := <-resultChannel
	log.Println("SyncContact", result)

	return nil //close the stream back to the client
}

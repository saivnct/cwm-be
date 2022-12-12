package apphttp

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"sol.go/cwm/dao"
	"sol.go/cwm/proto/cwmSIPPb"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/pubsub"
	"time"
)

type TestHTTPController struct {
}

func (sv *TestHTTPController) SendMsg(ctx *gin.Context) {
	name := ctx.Param("name")
	fmt.Printf("Http SendMsg function was invoked with %v\n", name)

	//region redis publish JSON
	//updateUserBalance := form.UpdateUserBalance{
	//	Username:   name,
	//	AmountUsdt: "10000000",
	//}
	//jsonString, err := json.Marshal(updateUserBalance)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status": "failed",
	//		"error":  err,
	//	})
	//	return
	//}
	//core.GetCache().RedisClient.Publish(context.Background(), static.EVENT_UPDATEUSERBALANCE, string(jsonString)).Err()
	//endregion

	//region redis publish proto
	cwmRequestHeader := cwmSIPPb.CWMRequestHeader{
		Method: cwmSIPPb.REQUEST_METHOD_METHOD_MESSAGE,
		From:   name,
		To:     "B",
	}

	msg := []byte("Hello World :)))))) \n Sáng 23/6, Tổng bí thư Nguyễn Phú Trọng cùng đại biểu Quốc hội thuộc đơn vị bầu cử số 1 TP Hà Nội tiếp xúc cử tri các quận Đống Đa, Ba Đình, Hai Bà Trưng. Về việc Trung ương khai trừ Đảng với Bộ trưởng Y tế Nguyễn Thanh Long và Chủ tịch UBND TP Hà Nội Chu Ngọc Anh, Tổng bí thư cho biết có ý kiến băn khoăn \"kỷ luật cách chức thì ai làm\" và cho rằng phải \"tìm người thay\" ")

	checksum := fmt.Sprintf("%x", md5.Sum(msg))
	log.Println("msg checksum", checksum)
	imType := cwmSignalMsgPb.SIGNAL_IM_TYPE_IM

	signalMessage := cwmSignalMsgPb.SignalMessage{
		MsgId:      "msg id",
		ThreadType: cwmSignalMsgPb.SIGNAL_THREAD_TYPE_SOLO,
		ImType:     imType,
		MsgDate:    time.Now().UnixMilli(),
		Checksum:   checksum,
		Data:       msg,
	}

	signalMessageData, err := proto.Marshal(&signalMessage)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
		return
	}

	cwmRequest := cwmSIPPb.CWMRequest{
		Header:  &cwmRequestHeader,
		Content: signalMessageData,
	}

	cwmRequestData, err := proto.Marshal(&cwmRequest)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  err,
		})
		return
	}

	//log.Println("base64", base64.StdEncoding.EncodeToString(cwmRequestData))
	//log.Println("hex", hex.EncodeToString(cwmRequestData))

	//https://stackoverflow.com/questions/3183841/base64-vs-hex-for-sending-binary-content-over-the-internet-in-xml-doc#:~:text=The%20difference%20between%20Base64%20and,it's%20more%20efficient%20than%20hex.
	dao.GetCache().RedisClient.Publish(
		context.Background(),
		pubsub.EVENT_CWMMSG,
		base64.StdEncoding.EncodeToString(cwmRequestData),
	).Err()
	//endregion

	//region send FCM
	//appfirebase.GetAppFireBase().SendMsg(name)
	//endregion

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (sv *TestHTTPController) Test(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

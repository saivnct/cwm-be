package pubsub

import (
	"context"
	"fmt"
	"log"
	"sol.go/cwm/dao"
)

const (
	EVENT_CWMMSG = "CWMMsg"
)

var Events = []string{EVENT_CWMMSG}

func StartSubscribe() {
	cache := dao.GetCache()
	subscriber := cache.RedisClient.Subscribe(context.Background(), Events...)
	userPubSubServer := UserPubSubHandler{}

	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				log.Println("redis Subscribe error", err)
			}

			switch msg.Channel {
			case EVENT_CWMMSG:
				fmt.Println("redis: received - EVENT_SIGNALMSG")

				userPubSubServer.HandleCWMMsg(msg.Payload)
				break
			default:
				fmt.Println("redis: unknown message", msg)
				break
			}

		}
	}()
}

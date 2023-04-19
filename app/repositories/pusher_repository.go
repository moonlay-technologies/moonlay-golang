package repositories

import (
	"encoding/json"
	"fmt"
	"order-service/app/models"
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

type PusherRepositoryInterface interface {
	Publish(*models.Pusher) error
}

type pusherRepository struct {
	pusherClient pusher.Client
}

func InitPusherRepository() PusherRepositoryInterface {
	return &pusherRepository{
		pusherClient: pusher.Client{
			AppID:   os.Getenv("PUSHER_APP_ID"),
			Key:     os.Getenv("PUSHER_KEY"),
			Secret:  os.Getenv("PUSHER_SECRET"),
			Cluster: os.Getenv("PUSHER_CLUSTER"),
			Secure:  true,
		},
	}
}

func (r *pusherRepository) Publish(request *models.Pusher) error {
	request.Channel = os.Getenv("PUSHER_CHANNEL")
	sJson, _ := json.Marshal(request)
	var data map[string]interface{}
	json.Unmarshal(sJson, &data)
	err := r.pusherClient.Trigger(request.Channel, os.Getenv("PUSHER_EVENT"), data)
	if err != nil {
		fmt.Println("pusher error = ", err.Error())
		return err
	}
	return nil
}

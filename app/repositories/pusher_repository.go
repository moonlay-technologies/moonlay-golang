package repositories

import (
	"fmt"
	"order-service/app/models/constants"
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

type PusherRepositoryInterface interface {
	Pubish(data map[string]string) error
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

func (r *pusherRepository) Pubish(data map[string]string) error {
	err := r.pusherClient.Trigger(constants.S3_EXPORT_CHANNEL, constants.S3_EXPORT_EVENT, data)
	if err != nil {
		fmt.Println("pusher error = ", err.Error())
		return err
	}
	return nil
}

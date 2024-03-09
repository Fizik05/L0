package nats

import (
	"encoding/json"

	"github.com/Fizik05/L0"
	"github.com/Fizik05/L0/pkg/repository"
	"github.com/Fizik05/L0/pkg/service"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

func NewSubscribeToChannel(clusterId, clientId, channelName string, repo *repository.Repository, service *service.Service) (stan.Conn, error) {
	natsURL := "localhost:4222"
	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(natsURL))
	if err != nil {
		logrus.Errorf("Error during connection to NATS-Streaming: %s", err.Error())
		return nil, err
	}

	msgHandler := func(msg *stan.Msg) {
		var order L0.Order

		if err := json.Unmarshal(msg.Data, &order); err != nil {
			logrus.Errorf("Error during unparsing json: %s", err.Error())
			return
		}

		if err := repo.SaveOrder(order); err != nil {
			logrus.Errorf("Error during posting in database: %s", err.Error())
			return
		}
		service.Cache.AddOrder(order.Order_uid, order)
		logrus.Println("New message")
	}

	_, err = sc.Subscribe(channelName, msgHandler)
	if err != nil {
		logrus.Errorf("Error during subscribing to channel: %s", err.Error())
		return nil, err
	}
	logrus.Printf("Succesful connect")
	return sc, nil
}

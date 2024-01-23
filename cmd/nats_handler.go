// nats_handler.go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"L0/entity"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
)

type NatsHandler struct {
	db *gorm.DB
}

func NewNatsHandler(db *gorm.DB) *NatsHandler {
	return &NatsHandler{db: db}
}

func (nh *NatsHandler) ConnectAndSubscribe() {
	nc, err := nats.Connect("nats://nats-streaming-l0:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("your_channel_name", nh.handleMessage)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to NATS Streaming. Waiting for messages on channel 'your_channel_name'...")

	select {}
}

func (nh *NatsHandler) handleMessage(msg *nats.Msg) {
	var order entity.Order

	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	err = nh.db.Save(&order).Error
	fmt.Println(&order)
	if err != nil {
		log.Println("Error saving to database:", err)
		return
	}
	err = nh.db.Save(&order.Payment).Error
	fmt.Println(&order.Payment)
	if err != nil {
		log.Println("Error saving to database:", err)
		return
	}
	err = nh.db.Save(&order.Delivery).Error
	fmt.Println(&order.Delivery)
	if err != nil {
		log.Println("Error saving to database:", err)
		return
	}
	for i := 0; i < len(order.Items); i++ {
		err = nh.db.Save(&order.Items[i]).Error
		if err != nil {
			log.Println("Error saving to database:", err)
			return
		}
	}

	log.Printf("Received and saved order with ID %d", order.OrderUid)
}

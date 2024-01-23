package main

import (
	"encoding/json"
	"log"
	"sync"

	"L0/entity"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
)

type NatsHandler struct {
	db         *gorm.DB
	orderCache *map[string]entity.Order
	cacheMutex *sync.RWMutex
}

func NewNatsHandler(db *gorm.DB, orderCache *map[string]entity.Order, cacheMutex *sync.RWMutex) *NatsHandler {
	return &NatsHandler{db: db, orderCache: orderCache, cacheMutex: cacheMutex}
}

func (nh *NatsHandler) ConnectAndSubscribe() {
	nc, err := nats.Connect("nats://nats-streaming-l0:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("your_channel_name", nh.handleMessage)
	if err != nil {
		log.Fatal("Failed to subscribe to channel:", err)
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

	if err := nh.saveToDatabase(&order); err != nil {
		log.Println("Error saving to database:", err)
		return
	}

	nh.cacheMutex.Lock()
	(*nh.orderCache)[order.OrderUid] = order
	nh.cacheMutex.Unlock()

	log.Printf("Received and saved order with ID %d", order.OrderUid)
}

func (nh *NatsHandler) saveToDatabase(order *entity.Order) error {
	tx := nh.db.Begin()

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&order.Payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&order.Delivery).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range order.Items {
		if err := tx.Save(&order.Items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

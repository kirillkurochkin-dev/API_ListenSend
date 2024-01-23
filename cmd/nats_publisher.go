package main

import (
	"L0/entity"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"time"
)

// PublishOrder отправляет сообщение с заказом в канал nats-streaming
func PublishOrder(orderData []byte) {
	// Подключение к nats-streaming
	nc, err := nats.Connect("nats://nats-streaming-l0:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Определение имени канала
	channelName := "your_channel_name"

	err = nc.Publish(channelName, orderData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Order published successfully to channel", channelName)
}

func PublishStart() {
	log.Print("Message delivered")
	order := entity.NewOrder()
	order.OrderUid = "b563feb7b2b84b6test"
	order.TrackNumber = "WBILMTESTTRACK"
	order.Entry = "WBIL"
	order.Delivery = entity.Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	}
	order.Payment = entity.Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	}
	order.Items = []entity.Item{
		{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	}
	order.Locale = "en"
	order.InternalSignature = ""
	order.CustomerId = "test@gmail.com"
	order.DeliveryService = "meest"
	order.ShardKey = "9"
	order.SmId = 99
	order.DateCreated, _ = time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")
	order.OffShard = "1"

	var orderBytes []byte

	for i := 1; i <= 10; i++ {
		if i%4 == 0 {
			orderBytes = []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x47, 0x6F, 0x2D, 0x4C, 0x61, 0x6E, 0x67}
		} else if i%3 == 0 {
			orderBytes, _ = json.Marshal([]byte{'t'})
		} else {

			order.CustomerId = "test" + strconv.Itoa(i) + "@gmail.com"
			order.Delivery.Email = "test" + strconv.Itoa(i) + "@gmail.com"

			order.OrderUid = "b" + strconv.Itoa(i) + "63feb7b2b84b6test"
			order.Payment.Transaction = "b" + strconv.Itoa(i) + "63feb7b2b84b6test"

			order.TrackNumber = "WBILMTESTTRACK" + strconv.Itoa(i)
			for j := 0; j < len(order.Items); j++ {
				order.Items[j].TrackNumber = "WBILMTESTTRACK" + strconv.Itoa(i)
			}

			orderBytes, _ = json.Marshal(order)
		}
		PublishOrder(orderBytes)
	}

}

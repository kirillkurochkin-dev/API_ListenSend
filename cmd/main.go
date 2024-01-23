package main

import (
	"L0/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"sync"
)

var db *gorm.DB
var err error

var cacheMutex sync.RWMutex
var orderCache map[string]entity.Order

func main() {
	db, err = gorm.Open("postgres", "host=postgres-L0 port=5432 user=postgres dbname=l0 sslmode=disable password=postgres")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&entity.Order{}, &entity.Delivery{}, &entity.Item{}, &entity.Payment{})

	//
	restoreCacheFromDB()
	//

	natsHandler := NewNatsHandler(db)

	go natsHandler.ConnectAndSubscribe()
	router := gin.Default()

	router.GET("/orders", getOrdersIds)
	router.POST("/orders", start)
	router.GET("/orders/:id", getOrderById)
	router.DELETE("/orders", deleteAllRecords)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}

func restoreCacheFromDB() {
	var orders []entity.Order
	var ordersCopy []entity.Order

	if err := db.Find(&orders).Error; err != nil {
		log.Println("Failed to retrieve orders from DB:", err)
		return
	}

	var delivery entity.Delivery
	var payment entity.Payment
	var items []entity.Item

	for _, order := range orders {
		fmt.Println(order.OrderUid)

		db.Where("email = ?", order.CustomerId).First(&delivery)
		db.Where("transaction = ?", order.OrderUid).First(&payment)
		db.Where("track_number = ?", order.TrackNumber).Find(&items)

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		ordersCopy = append(ordersCopy, order)
	}

	cacheMutex.Lock()
	orderCache = make(map[string]entity.Order)
	for _, order := range ordersCopy {
		orderCache[order.OrderUid] = order
	}
	cacheMutex.Unlock()

	log.Println("Cache restored from DB")
}

func start(c *gin.Context) {
	PublishStart()
}

func getOrderById(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var order entity.Order
	var delivery entity.Delivery
	var payment entity.Payment
	var items []entity.Item

	db.Where("order_uid = ?", orderID).First(&order)
	db.Where("email = ?", order.CustomerId).First(&delivery)
	db.Where("transaction = ?", orderID).First(&payment)
	db.Where("track_number = ?", order.TrackNumber).Find(&items)

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	fmt.Println(order)

	c.JSON(http.StatusOK, order)
}

func getOrdersIds(c *gin.Context) {
	var orders []entity.Order

	db.Find(&orders)

	var ids []string

	for _, order := range orders {
		ids = append(ids, order.OrderUid)
	}

	c.JSON(http.StatusOK, ids)
}

func deleteAllRecords(c *gin.Context) {
	if err := db.Delete(entity.Order{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete records"})
		return
	}
	if err := db.Delete(entity.Delivery{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete records"})
		return
	}
	if err := db.Delete(entity.Item{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete records"})
		return
	}
	if err := db.Delete(entity.Payment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All records deleted successfully"})
}

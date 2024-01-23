package main

import (
	"L0/entity"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db         *gorm.DB
	err        error
	cacheMutex sync.RWMutex
	orderCache map[string]entity.Order
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	initDB()
	defer db.Close()

	db.AutoMigrate(&entity.Order{}, &entity.Delivery{}, &entity.Item{}, &entity.Payment{})

	restoreCacheFromDB()

	natsHandler := NewNatsHandler(db, &orderCache, &cacheMutex)
	go natsHandler.ConnectAndSubscribe()

	router := gin.Default()
	router.GET("api/orders", getOrdersIds)
	router.POST("api/orders", start)
	router.GET("api/orders/:id", getOrderById)
	router.DELETE("api/orders", deleteAllRecords)

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/order/:id", getOrderPage)
	router.GET("/orders", getOrdersPage)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getOrderPage(c *gin.Context) {
	orderID := c.Param("id")

	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	order, found := orderCache[orderID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.HTML(http.StatusOK, "order.html", gin.H{"Order": order})
}

func getOrdersPage(c *gin.Context) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	var orders []entity.Order
	for _, order := range orderCache {
		orders = append(orders, order)
	}

	c.HTML(http.StatusOK, "orders.html", gin.H{"Orders": orders})
}

func initDB() {
	db, err = gorm.Open("postgres", "host=postgres-L0 port=5432 user=postgres dbname=l0 sslmode=disable password=postgres")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}

func restoreCacheFromDB() {
	var orders []entity.Order
	var ordersCopy []entity.Order

	if err := db.Find(&orders).Error; err != nil {
		log.Println("Failed to retrieve orders from DB:", err)
		return
	}

	for _, order := range orders {
		var delivery entity.Delivery
		var payment entity.Payment
		var items []entity.Item

		db.Where("email = ?", order.CustomerId).First(&delivery)
		db.Where("transaction = ?", order.OrderUid).First(&payment)
		db.Where("track_number = ?", order.TrackNumber).Find(&items)

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		ordersCopy = append(ordersCopy, order)
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	orderCache = make(map[string]entity.Order)
	for _, order := range ordersCopy {
		orderCache[order.OrderUid] = order
	}

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

	cacheMutex.RLock()
	order, found := orderCache[orderID]
	cacheMutex.RUnlock()

	if !found {
		var orderFromDB entity.Order
		var delivery entity.Delivery
		var payment entity.Payment
		var items []entity.Item

		if err := db.Where("order_uid = ?", orderID).First(&orderFromDB).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		db.Where("email = ?", order.CustomerId).First(&delivery)
		db.Where("transaction = ?", orderID).First(&payment)
		db.Where("track_number = ?", order.TrackNumber).Find(&items)

		orderFromDB.Delivery = delivery
		orderFromDB.Payment = payment
		orderFromDB.Items = items

		cacheMutex.Lock()
		defer cacheMutex.Unlock()
		orderCache[orderFromDB.OrderUid] = orderFromDB
		order = orderFromDB
	}

	c.JSON(http.StatusOK, order)
}

func getOrdersIds(c *gin.Context) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	var ids []string
	for _, order := range orderCache {
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

	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	orderCache = make(map[string]entity.Order)

	c.JSON(http.StatusOK, gin.H{"message": "All records deleted successfully"})
}

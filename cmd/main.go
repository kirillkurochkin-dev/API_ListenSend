package main

import (
	"L0/entity"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

var db *gorm.DB
var err error

func main() {
	// Подключение к базе данных
	db, err = gorm.Open("postgres", "host=postgres-L0 port=5432 user=postgres dbname=l0 sslmode=disable password=postgres")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы пользователей
	db.AutoMigrate(&entity.Order{}, &entity.Delivery{}, &entity.Item{}, &entity.Payment{})

	// Создание экземпляра Gin
	router := gin.Default()

	// Обработчик для получения списка заказов
	router.GET("/orders", getOrders)

	// Запуск сервера
	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getOrders(c *gin.Context) {
	var orders []entity.Order
	db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

# Задание L0

## Технологический стек
- Go 1.21
- NATs Streaming - реализация системы очередей сообщений
- Gin - веб-фреймворк, альтернатива стандартному пакету net/http
- GORM - ORM фреймворк для Go
- PostgresSQL - реляционные базы данных
- Docker, Docker Compose - контейнеризация сервисов приложения

## Установка и запуск проекта

Необходимо настроенная система виртуализации, установленный Docker Desktop(скачать и установить можно с официального сайта https://www.docker.com/products/docker-desktop/)

1. Клонируйте репозиторий проекта на свою локальную машину.
2. Запустите коммандную строку и перейдите в коррень директории с проектом.
3. Введите следующую команду, которая подготовит и запустит приложение на вашей локальной машине

   ```
   $  docker-compose up
   ```
4. Приложение будет запущено на порту 8080 и готово принимать http-запросы.
5. Спсиок доступных эндпоинтов предоставлен ниже.

router := gin.Default()
router.GET("api/orders", getOrdersIds)
router.POST("api/orders", start)
router.GET("api/orders/:id", getOrderById)
router.DELETE("api/orders", deleteAllRecords)

	router.LoadHTMLGlob("templates/*.html")

router.GET("/order/:id", getOrderPage)
router.GET("/orders", getOrdersPage)

### Эндпоинты API
---
- GET    /api/orders - Получение данных обо всех ID заказов
- POST   /api/orders - Запуск процесса отправки и принятия в сообщений Nats-Streaming (отправляет как корректные, так и не корректные сообщения)
- GET    /api/orders/:id - Получение данных об одном заказе по id
- DELETE /car/{id} - Удаление данных обо всех заказах
---
### Эндпоинты сервиса
---
- GET    /order/:id - Веб страница конкретного заказа
- POST   /orders - Веб страница отображения всех заказов
---

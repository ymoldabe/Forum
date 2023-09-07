package main

import (
	"log"

	"git/ymoldabe/forum/configs"
	"git/ymoldabe/forum/internal/handler"
	"git/ymoldabe/forum/internal/server"
	"git/ymoldabe/forum/internal/service"
	"git/ymoldabe/forum/internal/store"
)

func main() {
	// Загрузка конфигурации приложения.
	config, err := configs.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Создание и инициализация базы данных.
	db, err := store.NewSqlite3(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание экземпляра хранилища данных.
	store := store.New(db)

	// Создание сервиса, который будет управлять бизнес-логикой.
	service := service.New(store)

	// Создание обработчика HTTP-запросов, связанного с бизнес-логикой.
	handler := handler.New(service)

	// Создание экземпляра сервера.
	srv := new(server.Server)

	// Запуск сервера, передавая порт и зарегистрированные маршруты.
	if err := srv.Run(config.Port, handler.InitRouters()); err != nil {
		log.Fatalf("Error in main: %s", err)
	}
}

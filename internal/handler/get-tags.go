package handler

import (
	"git/ymoldabe/forum/internal/service"
	"git/ymoldabe/forum/models"
	"sync"
)

// userUrlBefore хранит предыдущие URL для каждого пользователя.
var (
	userUrlBefore = make(map[string]string)
	mu            sync.Mutex
)

// Handler представляет обработчик HTTP-запросов.
type Handler struct {
	service *service.Service
}

// New создает и возвращает новый экземпляр обработчика с заданным сервисом.
func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// getTags выполняет поиск постов, содержащих заданные теги, в переданных данных.
func getTags(tags []string, form []models.GetAllPosts) []models.GetAllPosts {
	foundData := []models.GetAllPosts{}

	searchMap := make(map[string]bool)

	for _, tag := range tags {
		searchMap[tag] = true
	}

	for _, post := range form {
		for _, postTag := range post.Tags {
			if searchMap[postTag] {
				foundData = append(foundData, post)
				break
			}
		}
	}
	return foundData
}

package server

import (
	"git/ymoldabe/forum/pkg"
	"net/http"
	"time"
)

// Server представляет структуру для управления HTTP-сервером.
type Server struct {
	httpServer *http.Server
}

// Run запускает HTTP-сервер на указанном порту с заданным обработчиком.
func (s *Server) Run(port string, handler http.Handler) error {
	// Создаем экземпляр HTTP-сервера с заданными параметрами.
	s.httpServer = &http.Server{
		Addr:           port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // Максимальный размер заголовков в байтах.
		ReadTimeout:    10 * time.Second, // Тайм-аут чтения запроса.
		WriteTimeout:   10 * time.Second, // Тайм-аут записи ответа.
	}

	// Выводим информацию о запуске сервера в лог.
	pkg.InfoLog.Printf("Server run on http://localhost%s", port)

	// Запускаем сервер и возвращаем ошибку (если есть).
	return s.httpServer.ListenAndServe()
}

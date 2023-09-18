package handler

import "net/http"

// InitRouters инициализирует роутеры и маршруты для обработки HTTP-запросов.
func (h *Handler) InitRouters() http.Handler {
	// Создаем новый маршрутизатор.
	mux := http.NewServeMux()

	// Создаем сервер для статических файлов и обслуживаем файлы из директории "./ui/static/".
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Создаем middleware цепочку для аутентификации.
	dynamic := aliceNew(h.requireAuthentication)

	// Добавляем обработчики для различных маршрутов.
	mux.HandleFunc("/", h.home)
	mux.Handle("/my-posts", dynamic.ThenFunc(h.MyPosts))
	mux.Handle("/my-likes", dynamic.ThenFunc(h.MyLikes))

	mux.HandleFunc("/post/view", h.postView)
	mux.Handle("/post/create", dynamic.ThenFunc(h.postCreate))
	mux.Handle("/post/coment-post", dynamic.ThenFunc(h.createPostComment))

	mux.Handle("/post/reaction", dynamic.ThenFunc(h.reaction))
	mux.Handle("/post/reaction-comment", dynamic.ThenFunc(h.reactionComment))

	mux.HandleFunc("/auth/sign-up", h.signUp)
	mux.HandleFunc("/auth/sign-in", h.signIn)
	mux.HandleFunc("/auth/google/login", h.googleLogin)
	mux.HandleFunc("/auth/google/callback", h.googleCallback)
	mux.HandleFunc("/auth/github/login", h.githubLogin)
	mux.HandleFunc("/auth/github/callback", h.githubCallback)

	mux.Handle("/logout", dynamic.ThenFunc(h.logout))

	// Создаем middleware цепочку для обработки ошибок и логирования.
	standard := aliceNew(h.recoverPanic, h.logRequest)

	// Возвращаем готовый HTTP-обработчик, который объединяет middleware и маршруты.
	return standard.Then(mux)
}

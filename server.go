package level0

import (
    "net/http"
    "time"
    "context"
)

// Стуктура для запуска http сервера
type Server struct {
    httpServer *http.Server
}

// Метод отвечающий за запуск работы сервера
func (s *Server) Run (port string) error {
    s.httpServer = &http.Server{
        Addr: ":" + port,
        MaxHeaderBytes: 1 << 20,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Запускаем бесконечный цикл для прослушивание всех входящих запросов
    // и их последующей обработки
    return s.httpServer.ListenAndServe()
}

// Метод отвечающий за остановку работы сервера
func (s *Server) Shutdown(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}

package main

import (
    "github.com/Wugdich/wugdsrv"
    "log"
)

func main() {
    // Инициализируем экземпляр сервера
    srv := new(wugdsrv.Server)
    
    // Запускаем сервер, если при запуске не возникло ошибки
    if err := srv.Run("8000"); err != nil {
        log.Fatalf("error occured while running the server: %s", err.Error())
}

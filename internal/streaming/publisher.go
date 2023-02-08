package streaming

import (
	"encoding"
	"encoding/json"
	"log"
	"wugdsrv/internal/pstgrs"

	stan "github.com/nats-io/stan.go"
)

type Publisher struct {
    sc *stan.Conn
    name string
}

// Функция для инициализации экземпляра структуры
func NewPublisher(s *stan.Conn) *Publisher {
    return &Publisher{
        name: "Publisher",
        sc: s,
    }
}

// Функция для публикации информации о заказах в базу данных
func (p *Publisher) Publish() {
    // Создаем переменные для структур,которые будем отправлять в таблицы
    order := pstgrs.Order{}
    // Заворачиваем в json
    orderData, err := json.Marshal(order)
    if err != nil {
        log.Printf("%s: json.Marshal error: %v\n", p.name, err)
    }
   
    // An asynchronous publish API
    ackHandler := func(ackedNuid string, err error) {
        if err != nil {
            log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
        } else {
            log.Printf("Received ack for msg id %s\n", ackedNuid)
        }
    }

    // Асинхронная публикация
    log.Printf("%s: publishing data ...\n", p.name)
    nuid, err := (*p.sc).PublishAsync("testChannel", orderData, ackHandler)
    if err != nil {
        log.Printf("%s: error publishing msg %s: %v\n", p.name, nuid, err.Error())
    }

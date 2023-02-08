package streaming

import (
    "log"
    "encoding/json"
    "time"
    "wugdsrv/internal/pstgrs"

    stan "github.com/nats-io/stan.go"
)

// Структура подписки
type Subscriber struct {
    sub      stan.Subscription // Интерфейс подписки
    db       *pstgrs.DB
    sc       *stan.Conn //Stan Connction - интерфейс подключения
    name     string
}

// Инициализация структуры подписки
func NewSubscriber(conn *stan.Conn) *Subscriber {
    return &Subscriber{
        name:   "Subscriber",
        sc:     conn,
    }
}

// Функция обработки входящего сообщения
func (s *Subscriber) messageHandler(data []byte) bool {
    curOrder := "order from db"
    err := json.Unmarshal(data, &curOrder)
    // Если не удалось декодировать
    if err != nil {
        log.Printf("%s: messageHandler() error, %v\n", s.name, err)
        return true // Все равно сообщаем о полученном собщении
    }
    log.Printf("%s: unmarshal Order to struct: %v\n", s.name, curOrder)

    _, err = s.db.AddOrder(curOrder)
    // Проверяем, что заказ корректно добавился в базу
    if err != nil {
        log.Printf("%s: unable to add order: %v\n", s.name, err)
        return false
    }
    return true
}

func (s *Subscriber) Subscribe() {
    var err error

    // Подписка к кластеру. Указываем сабжект и опции.
    s.sub, err = (*s.sc).Subscribe("testChannel", func(m *stan.Msg) {
        log.Printf("%s: recieved a message!\n", s.name)
        if s.messageHandler(m.Data) {
            err := m.Ack() // Функция подтверждает сообщение
            if err != nil {
                log.Printf("%s ack() err: %s", s.name, err)
            }
        }
    },
    // Задаем параметры подписки
    stan.AckWait(30*time.Second), // Таймаут ожидания ack()
    stan.DurableName("natsDurable"), // Долговечные подписки позволяют достовлять только подтвержденные сообщения
    stan.SetManualAckMode(), // Позволяет клиенту контролировать его ack() для доставленных сообщений
    stan.MaxInflight(5))    // Устанавилвает максимальное кол-во сообщений, которые могут ожидать подтверждения

    if err != nil {
        log.Printf("%s: error: %v\n", s.name, err)
    }
    log.Printf("%s: subscribed to subject %s\n", s.name, "testSubject")
}

func (s *Subscriber) Unsubscribe() {
    if s.sub != nil {
        s.sub.Unsubscribe() // remove interests on the subscription
    }
}

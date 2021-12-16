package amqp

import (
	"fmt"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func onconnect(conn *amqp.Connection) {
	fmt.Printf("connect callback\n")
}

func onConnectionLost(conn *amqp.Connection, err *amqp.Error) {
	fmt.Printf("connection lost callback\n")
}

func TestAmqp(t *testing.T) {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	amqpw := NewClient(NewClientOptions().
		SetAddr(url).
		SetOnConnect(onconnect).SetOnConnectionLost(onConnectionLost))

	defer amqpw.Close()
	count := 0
	for {
		time.Sleep(3 * time.Second)
		conn, err := amqpw.Conn()
		if err != nil {
			t.Error(err)
			continue
		} else {
			t.Log("amqp ok")
		}
		ch, err := conn.Channel()
		if err != nil {
			t.Error(err)
			continue
		}

		ch.Publish("", "hello", false, false, amqp.Publishing{})
		ch.Close()
		count++
		if count > 100 {
			break
		}

	}

}

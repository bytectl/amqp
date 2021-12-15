package amqp

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestAmqp(t *testing.T) {

	amqpw, err := New("amqp://guest:guest@localhost:5672", amqp.Config{})
	if err != nil {
		t.Error(err)
		return
	}
	defer amqpw.Close()
	count := 0
	for {
		time.Sleep(1 * time.Second)
		conn, err := amqpw.Conn()
		if err != nil {
			t.Error(err)
			continue
		} else {
			t.Log("amqp ok11")
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

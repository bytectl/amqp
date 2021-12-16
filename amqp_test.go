package amqp

import (
	"fmt"
	"os"
	"testing"
	"time"

	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func reconnectCallback(ctx context.Context, conn *amqp.Connection) {
	fmt.Printf("reconnect callback\n")
	go func() {
		select {
		case <-ctx.Done():
			fmt.Printf("reconnect callback done\n")
			return
		case <-time.After(100 * time.Second):
			fmt.Printf("reconnect callback timeout\n")
		}
	}()
}

func TestAmqp(t *testing.T) {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	amqpw, err := New(url, amqp.Config{}, reconnectCallback)
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

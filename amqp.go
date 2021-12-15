package amqp

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
)

// HandlerFunc amqp message handler function.
type HandlerFunc func(ctx context.Context, delivery *amqp.Delivery)

// Amqp
type Amqpw struct {
	connection      *amqp.Connection
	isReady         bool
	notifyConnClose chan *amqp.Error
	done            chan bool
}

// New create amqp connection
func New(addr string) (*Amqpw, error) {
	connection, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	amqpw := &Amqpw{
		connection:      connection,
		isReady:         true,
		notifyConnClose: make(chan *amqp.Error),
		done:            make(chan bool),
	}
	connection.NotifyClose(amqpw.notifyConnClose)
	go amqpw.handleReconnect(addr)
	return amqpw, nil
}

// Close close amqp connection
func (amqpw *Amqpw) Close() error {
	if !amqpw.isReady {
		return errAlreadyClosed
	}
	err := amqpw.connection.Close()
	if err != nil {
		return err
	}
	amqpw.isReady = false
	close(amqpw.done)
	amqpw.done = nil
	amqpw.connection = nil
	return nil
}

// handleReconnect 等待notifyConnClose通知, 并尝试重连
func (amqpw *Amqpw) handleReconnect(addr string) {
	for {
		select {

		case <-amqpw.done:
			// amqpw 退出
			break
		case <-amqpw.notifyConnClose:
			// 连接出错
			fmt.Println("Connection closed. Reconnecting...")
			amqpw.isReady = false
			amqpw.connection = nil
		case <-time.After(reconnectDelay):
			// 重连
			if amqpw.isReady {
				continue
			}
			fmt.Println("Reconnecting...")

			connection, err := amqp.Dial(addr)
			if err != nil {
				fmt.Println("Failed to reconnect: %v", err)
				return
			}
			amqpw.connection = connection
			amqpw.isReady = true
		}
	}
}

// Conn get amqp connection
func (amqpw *Amqpw) Conn() *amqp.Connection {
	if amqpw.isReady {
		return amqpw.connection
	}
	return nil
}

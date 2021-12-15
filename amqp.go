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
	reconnectDelay = 3 * time.Second
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
func New(url string, c amqp.Config) (*Amqpw, error) {
	connection, err := amqp.DialConfig(url, c)
	if err != nil {
		return nil, err
	}
	amqpw := &Amqpw{
		isReady:         true,
		notifyConnClose: make(chan *amqp.Error),
		done:            make(chan bool),
	}

	amqpw.changeConnection(connection)
	go amqpw.handleReconnect(url, c)
	return amqpw, nil
}

func (amqpw *Amqpw) changeConnection(connection *amqp.Connection) {
	amqpw.connection = connection
	amqpw.notifyConnClose = make(chan *amqp.Error)
	amqpw.connection.NotifyClose(amqpw.notifyConnClose)
}

// Close close amqp connection
func (amqpw *Amqpw) Close() error {
	if !amqpw.isReady {
		return errAlreadyClosed
	}
	amqpw.done <- true
	err := amqpw.connection.Close()
	if err != nil {
		return err
	}
	amqpw.isReady = false
	amqpw.connection = nil
	return nil
}

// handleReconnect 等待notifyConnClose通知, 并尝试重连
func (amqpw *Amqpw) handleReconnect(url string, c amqp.Config) {
	for {
		select {
		case <-amqpw.done:
			// amqpw 退出
			return
		case <-amqpw.notifyConnClose:
			// 连接出错
			fmt.Printf("Connection closed. Reconnecting...\n")
			amqpw.isReady = false
			amqpw.connection = nil
			connection, err := amqp.DialConfig(url, c)
			if err != nil {
				fmt.Printf("Failed to reconnect: %v\n", err)
				time.Sleep(reconnectDelay)
				continue
			}
			amqpw.changeConnection(connection)
			amqpw.isReady = true
		}
	}
}

// Conn get amqp connection
func (amqpw *Amqpw) Conn() (*amqp.Connection, error) {
	if amqpw.isReady {
		return amqpw.connection, nil
	}
	return nil, errNotConnected
}

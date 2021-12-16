package amqp

import (
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

// Client amqp client
type Client struct {
	clientOptions   *ClientOptions
	connection      *amqp.Connection
	isReady         bool
	notifyConnClose chan *amqp.Error
	done            chan bool
}

// New create amqp connection
func New(o *ClientOptions) *Client {
	client := &Client{
		isReady:         false,
		notifyConnClose: make(chan *amqp.Error),
		done:            make(chan bool),
		clientOptions:   o,
	}
	go client.handleReconnect()
	return client
}

func (client *Client) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error)
	client.connection.NotifyClose(client.notifyConnClose)
}

// Close close amqp connection
func (client *Client) Close() error {
	if !client.isReady {
		return errAlreadyClosed
	}
	client.done <- true
	err := client.connection.Close()
	if err != nil {
		return err
	}
	client.isReady = false
	client.connection = nil
	return nil
}

// handleReconnect 等待notifyConnClose通知, 并尝试重连
func (client *Client) handleReconnect() {
	onConnectionLost := client.clientOptions.OnConnectionLost
	onConnect := client.clientOptions.OnConnect
	for {
		if !client.isReady {
			c := amqp.Config{
				Vhost:      client.clientOptions.Vhost,
				ChannelMax: client.clientOptions.ChannelMax,
				FrameSize:  client.clientOptions.FrameSize,
				Heartbeat:  client.clientOptions.Heartbeat,
				Locale:     client.clientOptions.Locale,
			}
			connection, err := amqp.DialConfig(client.clientOptions.Addr, c)
			if err != nil {
				fmt.Printf("Failed to reconnect: %v\n", err)
				time.Sleep(reconnectDelay)
				continue
			}
			client.changeConnection(connection)
			client.isReady = true
			if onConnect != nil {
				onConnect(connection)
			}
		}

		select {
		case <-client.done:
			// client 退出
			if onConnectionLost != nil {
				onConnectionLost(client.connection, amqp.ErrClosed)
			}
			return
		case err := <-client.notifyConnClose:
			// 连接出错
			fmt.Printf("Connection closed. Reconnecting...\n")
			client.isReady = false
			client.connection = nil
			if onConnectionLost != nil {
				onConnectionLost(client.connection, err)
			}
		}

	}
}

// Conn get amqp connection
func (client *Client) Conn() (*amqp.Connection, error) {
	if client.isReady {
		return client.connection, nil
	}
	return nil, errNotConnected
}

// Channel get amqp channel
func (client *Client) Channel() (*amqp.Channel, error) {
	if client.isReady {
		return client.connection.Channel()
	}
	return nil, errNotConnected
}

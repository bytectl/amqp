package amqp

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OnConnectionLostHandle func(*amqp.Connection, *amqp.Error)
type OnConnectHandle func(*amqp.Connection)
type ClientOptions struct {
	Addr             string
	Vhost            string
	ChannelMax       int           // 0 max channels means 2^16 - 1
	FrameSize        int           // 0 max bytes means unlimited
	Heartbeat        time.Duration // less than 1s uses the server's interval
	Properties       amqp.Table
	Locale           string
	OnConnect        OnConnectHandle
	OnConnectionLost OnConnectionLostHandle
}

func NewClientOptions() *ClientOptions {
	o := &ClientOptions{
		Addr:             "amqp://guest:guest@localhost:5672/",
		Vhost:            "/",
		ChannelMax:       0,
		FrameSize:        0,
		Heartbeat:        1 * time.Second,
		OnConnect:        nil,
		OnConnectionLost: nil,
	}
	return o
}

//SetAddr set amqp addr
func (o *ClientOptions) SetAddr(addr string) *ClientOptions {
	o.Addr = addr
	return o
}

//SetVhost set amqp vhost
func (o *ClientOptions) SetVhost(vhost string) *ClientOptions {
	o.Vhost = vhost
	return o
}

//SetChannelMax set amqp channel max
func (o *ClientOptions) SetChannelMax(channelMax int) *ClientOptions {
	o.ChannelMax = channelMax
	return o
}

//SetFrameSize set amqp frame size
func (o *ClientOptions) SetFrameSize(frameSize int) *ClientOptions {
	o.FrameSize = frameSize
	return o
}

//SetHeartbeat set amqp heartbeat
func (o *ClientOptions) SetHeartbeat(heartbeat time.Duration) *ClientOptions {
	o.Heartbeat = heartbeat
	return o
}

//SetProperties set amqp properties
func (o *ClientOptions) SetProperties(properties amqp.Table) *ClientOptions {
	o.Properties = properties
	return o
}

//SetLocale set amqp locale
func (o *ClientOptions) SetLocale(locale string) *ClientOptions {
	o.Locale = locale
	return o
}

//SetOnConnect set amqp on connect
func (o *ClientOptions) SetOnConnect(onConnect OnConnectHandle) *ClientOptions {
	o.OnConnect = onConnect
	return o
}

//SetOnConnectionLost set amqp on connection lost
func (o *ClientOptions) SetOnConnectionLost(onConnectionLost OnConnectionLostHandle) *ClientOptions {
	o.OnConnectionLost = onConnectionLost
	return o
}

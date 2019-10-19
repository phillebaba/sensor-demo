package api

import (
	"github.com/golang/protobuf/ptypes/empty"
)

type Server struct {
	Broker *Broker
}

func (s *Server) GetTemperature(e *empty.Empty, stream TemperatureService_GetTemperatureServer) error {
	msgCh := s.Broker.Subscribe()
	var err error

	for {
		temperature := (<-msgCh).(float64)
		if err = stream.Send(&TemperatureResponse{Temperature: temperature}); err != nil {
			s.Broker.Unsubscribe(msgCh)
			break
		}
	}

	return err
}

type Broker struct {
	stopCh    chan struct{}
	publishCh chan interface{}
	subCh     chan chan interface{}
	unsubCh   chan chan interface{}
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan interface{}, 1),
		subCh:     make(chan chan interface{}, 1),
		unsubCh:   make(chan chan interface{}, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan interface{}]struct{}{}
	for {
		select {
		case <-b.stopCh:
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

func (b *Broker) Subscribe() chan interface{} {
	msgCh := make(chan interface{}, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan interface{}) {
	b.unsubCh <- msgCh
	close(msgCh)
}

func (b *Broker) Publish(msg interface{}) {
	b.publishCh <- msg
}

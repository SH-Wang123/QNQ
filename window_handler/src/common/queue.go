package common

import (
	"time"
)

type MessageQueue interface {
	//Send The message into the MessageQueue
	Send(message interface{})

	//Pull The message with the given size and timeout
	Pull(size int, timeout time.Duration) []interface{}

	//Size The current number of messages in MessageQueue
	Size() int

	// Capacity The maximum number if messages in MessageQueue
	Capacity() int
}

type QMessageQueue struct {
	queue    chan interface{}
	capacity int
}

func (m *QMessageQueue) Send(message interface{}) {
	select {
	case m.queue <- message:

	default:

	}
}

func (m *QMessageQueue) Pull(size int, timeout time.Duration) []interface{} {
	ret := make([]interface{}, 0)
	for i := 0; i < size; i++ {
		select {
		case msg := <-m.queue:
			ret = append(ret, msg)
		case <-time.After(timeout):
			return ret
		}
	}
	return ret
}

func (m *QMessageQueue) Size() int {
	return len(m.queue)
}

func (m *QMessageQueue) Capacity() int {
	return m.capacity
}

func NewMessageQueue(capacity int) MessageQueue {
	var mq MessageQueue
	mq = &QMessageQueue{
		queue:    make(chan interface{}, capacity),
		capacity: capacity,
	}
	return mq
}

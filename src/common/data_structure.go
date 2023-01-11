package common

import "errors"

type vector[T any] []T

const MaxLength = 200

type Queue[T any] struct {
	data   vector[T]
	head   int
	tail   int
	length int
}

func (queue *Queue[T]) InitQueue() {
	queue.data = make(vector[T], 200)
	queue.head = 0
	queue.tail = 0
	queue.length = 0
}

func (queue *Queue[T]) Enqueue(data T) error {
	if queue.length == MaxLength {
		return errors.New("Queue is full ")
	}
	queue.data[queue.tail] = data
	queue.tail = (queue.tail + 1) % MaxLength
	queue.length++
	return nil
}

func (queue *Queue[T]) Dequeue(data *T) error {
	if queue.length == 0 {
		return errors.New("Queue is empty ")
	}
	data = &(queue.data[queue.head])
	queue.head = (queue.head + 1) % MaxLength
	queue.length--
	return nil
}

package logexporter

import (
	"container/list"
	"sync"
)

type Queue interface {
	Enqueue(data string)
	Dequeue() string
	N() int
}

type queue struct {
	m    sync.Mutex
	list list.List
}

func (q *queue) Enqueue(data string) {
	q.m.Lock()
	defer q.m.Unlock()
	q.list.PushBack(data)
}

func (q *queue) Dequeue() string {
	q.m.Lock()
	defer q.m.Unlock()
	e := q.list.Front()
	data := e.Value.(string)
	q.list.Remove(e)
	return data
}

func (q *queue) N() int {
	q.m.Lock()
	defer q.m.Unlock()
	n := q.list.Len()
	return n
}

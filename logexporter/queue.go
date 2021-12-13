package main

import (
	"container/list"
	"sync"
)

type Queue interface {
	Enqueue(data string)
	Dequeue() string
	N() int
	End()
	IsEnded() bool
}

type queue struct {
	m       sync.Mutex
	list    list.List
	isEnded bool
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

func (q *queue) End() {
	q.m.Lock()
	defer q.m.Unlock()
	q.isEnded = true
}

func (q *queue) IsEnded() bool {
	q.m.Lock()
	defer q.m.Unlock()
	return q.isEnded
}

func NewQueue() Queue {
	return &queue{m: sync.Mutex{}, list: *list.New(), isEnded: false}
}
